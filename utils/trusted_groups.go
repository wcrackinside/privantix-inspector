package utils

import (
	"encoding/json"
	"os"
	"strings"
)

type TrustedGroup struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type TrustedGroupsFile struct {
	Groups []TrustedGroup `json:"groups"`
}

// LoadTrustedGroups loads trusted principals from a JSON file and returns a set (lowercase) for matching.
func LoadTrustedGroups(path string) (map[string]struct{}, error) {
	tg, err := readTrustedGroupsFile(path)
	if err != nil {
		return nil, err
	}
	set := make(map[string]struct{})
	for _, g := range tg.Groups {
		for _, m := range g.Members {
			m = strings.TrimSpace(m)
			if m == "" {
				continue
			}
			keyFull := strings.ToLower(m)
			set[keyFull] = struct{}{}
			if idx := strings.LastIndex(m, "\\"); idx != -1 && idx+1 < len(m) {
				keyUser := strings.ToLower(m[idx+1:])
				if keyUser != "" {
					set[keyUser] = struct{}{}
				}
			}
		}
	}
	return set, nil
}

func readTrustedGroupsFile(path string) (*TrustedGroupsFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tg TrustedGroupsFile
	if err := json.Unmarshal(b, &tg); err != nil {
		return nil, err
	}
	return &tg, nil
}

// ExtractPrincipal gets the user/group from ACL string "DOMAIN\user:(I)(F)" -> "DOMAIN\user"
func ExtractPrincipal(acl string) string {
	idx := strings.Index(acl, ":")
	if idx == -1 {
		return strings.TrimSpace(acl)
	}
	return strings.TrimSpace(acl[:idx])
}

// PrincipalToMember extracts the member part for matching: "DOMAIN\carce" -> "carce"
func PrincipalToMember(principal string) string {
	principal = strings.TrimSpace(principal)
	if idx := strings.LastIndex(principal, "\\"); idx != -1 && idx+1 < len(principal) {
		return principal[idx+1:]
	}
	return principal
}

// ClassifyPrincipals splits ACL principals into trusted and outside based on trustedSet.
func ClassifyPrincipals(acls []string, trustedSet map[string]struct{}) (trusted, outside []string) {
	seen := make(map[string]bool)
	for _, acl := range acls {
		principal := ExtractPrincipal(acl)
		if principal == "" {
			continue
		}
		if strings.Contains(principal, "archivos") || strings.Contains(principal, "processed") {
			continue
		}
		if seen[principal] {
			continue
		}
		seen[principal] = true
		principalNorm := strings.ToLower(principal)
		memberNorm := strings.ToLower(PrincipalToMember(principal))
		_, fullMatch := trustedSet[principalNorm]
		_, memberMatch := trustedSet[memberNorm]
		if fullMatch || memberMatch {
			trusted = append(trusted, principal)
		} else {
			outside = append(outside, principal)
		}
	}
	return trusted, outside
}

// ComplianceStatus returns "compliant" if no outside principals, else "non_compliant".
func ComplianceStatus(trusted, outside []string) string {
	if len(outside) == 0 {
		return "compliant"
	}
	return "non_compliant"
}
