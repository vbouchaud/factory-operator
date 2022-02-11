package ldap

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"reflect"
	"sort"
	"strings"
)

const (
	ScopeBaseObject   = "base"
	ScopeSingleLevel  = "single"
	ScopeWholeSubtree = "sub"
)

var scopeMap = map[string]int{
	ScopeBaseObject:   0,
	ScopeSingleLevel:  1,
	ScopeWholeSubtree: 2,
}

const (
	objectClass  = "objectClass"
	uniqueMember = "uniqueMember"
	description  = "description"

	groupClassValue = "groupOfUniqueNames"
)

const (
	errGroupNotFound     = "group name was not found"
	errGroupNotExtracted = "group name could not be extracted"
)

type Client struct {
	ldapURL               string
	bindDN                string
	bindPassword          string
	groupSearchBase       string
	groupSearchScope      string
	groupSearchFilter     string
	groupNameProperty     string
	groupSearchAttributes []string
}

func NewInstance(
	ldapURL,
	bindDN,
	bindPassword,
	groupSearchBase,
	groupSearchScope,
	groupSearchFilter,
	groupNameProperty string,
	groupSearchAttributes []string,
) *Client {
	s := &Client{
		ldapURL:               ldapURL,
		bindDN:                bindDN,
		bindPassword:          bindPassword,
		groupSearchBase:       groupSearchBase,
		groupSearchScope:      groupSearchScope,
		groupSearchFilter:     groupSearchFilter,
		groupNameProperty:     groupNameProperty,
		groupSearchAttributes: groupSearchAttributes,
	}

	return s
}

func (s *Client) bind() (*ldap.Conn, error) {
	l, err := ldap.DialURL(s.ldapURL)
	if err != nil {
		return nil, err
	}

	err = l.Bind(s.bindDN, s.bindPassword)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *Client) groupExists(name string) (bool, error, *ldap.Entry) {
	l, err := s.bind()
	if err != nil {
		return false, err, nil
	}

	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		s.groupSearchBase,
		scopeMap[s.groupSearchScope],
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf(s.groupSearchFilter, name),
		s.groupSearchAttributes,
		nil,
	)
	result, err := l.Search(searchRequest)
	if err != nil {
		return false, err, nil
	}

	if len(result.Entries) == 0 {
		return false, nil, nil
	} else if len(result.Entries) > 1 {
		return false, fmt.Errorf("too many entries returned"), nil
	}

	return true, nil, result.Entries[0]
}

func (s *Client) createGroup(groupDN, desc string, members []string) error {
	l, err := s.bind()
	if err != nil {
		return err
	}

	defer l.Close()

	addRequest := ldap.NewAddRequest(groupDN, nil)
	addRequest.Attribute(objectClass, []string{groupClassValue})
	addRequest.Attribute(description, []string{desc})
	addRequest.Attribute(uniqueMember, members)

	if err := l.Add(addRequest); err != nil {
		return err
	}

	return nil
}

func (s *Client) modifyGroup(groupDN, desc string, members []string) error {
	l, err := s.bind()
	if err != nil {
		return err
	}

	defer l.Close()

	modifyRequest := ldap.NewModifyRequest(groupDN, nil)
	modifyRequest.Replace(objectClass, []string{groupClassValue})
	modifyRequest.Replace(description, []string{desc})
	modifyRequest.Replace(uniqueMember, members)

	if err := l.Modify(modifyRequest); err != nil {
		return err
	}

	return nil
}

func (s *Client) deleteGroup(groupDN string) error {
	l, err := s.bind()
	if err != nil {
		return err
	}

	defer l.Close()

	delRequest := ldap.NewDelRequest(groupDN, nil)

	if err := l.Del(delRequest); err != nil {
		return err
	}

	return nil
}

func (s *Client) DeleteGroup(name string) error {
	groupDN := fmt.Sprintf("%s=%s,%s", s.groupNameProperty, name, s.groupSearchBase)

	exists, err, _ := s.groupExists(name)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(errGroupNotFound)
	}

	return s.deleteGroup(groupDN)
}

func sanitize(a []string) []string {
	var res []string

	for _, item := range a {
		res = append(res, strings.ToLower(item))
	}

	sort.Strings(res)

	return res
}

func (s *Client) ReconcileGroup(name, desc string, members []string) (string, error, bool) {
	groupDN := fmt.Sprintf("%s=%s,%s", s.groupNameProperty, name, s.groupSearchBase)

	exists, err, entry := s.groupExists(name)
	if err != nil {
		return groupDN, err, false
	}

	if exists {
		changed := entry.GetAttributeValue(description) != desc || !reflect.DeepEqual(sanitize(entry.GetAttributeValues(uniqueMember)), sanitize(members))
		if changed {
			return groupDN, s.modifyGroup(groupDN, desc, members), changed
		}
		return groupDN, nil, false
	}

	return groupDN, s.createGroup(groupDN, desc, members), true
}

func IsNotFound(err error) bool {
	return err.Error() == errGroupNotFound
}
