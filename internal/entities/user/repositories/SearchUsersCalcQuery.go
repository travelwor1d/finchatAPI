package repositories

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	srequest "github.com/finchatapp/finchat-api/internal/entities/_shared/models/request"
	baserepo "github.com/finchatapp/finchat-api/internal/entities/_shared/repositories"
)

func (s *repository) SearchUsersCalcQuery(ctx context.Context, userID int, filters []srequest.CustomFilterItem) (string, baserepo.SortDefinitionFunc, []string, error) {

	andConditions := make([]string, 0)
	//query
	q := fmt.Sprintf(`
	SELECT 
	u.id,
	u.is_active,
	u.first_name,
	u.last_name,
	u.phone_number,
	u.country_code,
	u.email,
	u.is_verified,
	u.user_type,
	u.profile_avatar,
	u.last_seen, 
	u.created_at, 
	u.updated_at,
	u.username,
	uc.id IS NOT NULL AS is_contact 

	FROM core.users u
	left join core.users_contacts uc on uc.user_id = %d AND u.id = uc.contact_id
	`, userID)
	params := make([]string, 0)

	andConditions = append(andConditions, fmt.Sprintf("(u.id <> :param_value_%d and u.deleted_at is null)", len(params)+1))
	params = append(params, strconv.Itoa(userID))
	//filters
	for _, filter := range filters {
		switch filter.Name {
		case "userTypes":
			{
				orCondition := make([]string, 0)
				for _, val := range filter.Values {
					orCondition = append(orCondition, fmt.Sprintf("(u.user_type = :param_value_%d)", len(params)+1))
					params = append(params, val)
				}
				if len(orCondition) > 0 {
					res := fmt.Sprintf("(%s)", strings.Join(orCondition, " or "))
					andConditions = append(andConditions, res)
				}
			}
		case "query":
			{
				if len(filter.Values) > 0 && filter.Values[0] != "" {
					val := strings.ToLower(filter.Values[0])
					orCondition := make([]string, 0)
					orCondition = append(orCondition, fmt.Sprintf("(lower(u.first_name) like  CONCAT('%%', :param_value_%d, '%%'))", len(params)+1))
					params = append(params, val)
					orCondition = append(orCondition, fmt.Sprintf("(lower(u.last_name) like  CONCAT('%%', :param_value_%d, '%%'))", len(params)+1))
					params = append(params, val)
					orCondition = append(orCondition, fmt.Sprintf("(lower(u.username) like  CONCAT('%%', :param_value_%d, '%%'))", len(params)+1))
					params = append(params, val)
					orCondition = append(orCondition, fmt.Sprintf("(lower(concat(u.first_name, ' ', u.last_name)) like  CONCAT('%%', :param_value_%d, '%%'))", len(params)+1))
					params = append(params, val)
					orCondition = append(orCondition, fmt.Sprintf("(lower(concat(u.last_name, ' ', u.first_name)) like  CONCAT('%%', :param_value_%d, '%%'))", len(params)+1))
					params = append(params, val)
					orCondition = append(orCondition, fmt.Sprintf("(lower(u.email) like  CONCAT('%%', :param_value_%d, '%%'))", len(params)+1))
					params = append(params, val)
					orCondition = append(orCondition, fmt.Sprintf("(lower(u.phone_number) like  CONCAT('%%', :param_value_%d, '%%'))", len(params)+1))
					params = append(params, val)

					if len(orCondition) > 0 {
						res := fmt.Sprintf("(%s)", strings.Join(orCondition, " or "))
						andConditions = append(andConditions, res)
					}
				}
			}
		case "ignoreContacts":
			{
				if len(filter.Values) > 0 && filter.Values[0] != "" {
					val := strings.ToLower(filter.Values[0])
					if val == "true" {
						orCondition := make([]string, 0)
						orCondition = append(orCondition, fmt.Sprintf("uc.id IS NULL"))
						params = append(params, val)
						if len(orCondition) > 0 {
							res := fmt.Sprintf("(%s)", strings.Join(orCondition, " or "))
							andConditions = append(andConditions, res)
						}
					}
				}
			}
		case "onlyContacts":
			{
				if len(filter.Values) > 0 && filter.Values[0] != "" {
					val := strings.ToLower(filter.Values[0])
					if val == "true" {
						orCondition := make([]string, 0)
						orCondition = append(orCondition, fmt.Sprintf("uc.id IS NOT NULL"))
						params = append(params, val)
						if len(orCondition) > 0 {
							res := fmt.Sprintf("(%s)", strings.Join(orCondition, " or "))
							andConditions = append(andConditions, res)
						}
					}
				}
			}
		case "contactId":
			{
				if len(filter.Values) > 0 && filter.Values[0] != "" {
					val := strings.ToLower(filter.Values[0])
					orCondition := make([]string, 0)
					orCondition = append(orCondition, fmt.Sprintf("u.id = :param_value_%d and uc.id IS NOT NULL", len(params)+1))
					params = append(params, val)
					if len(orCondition) > 0 {
						res := fmt.Sprintf("(%s)", strings.Join(orCondition, " or "))
						andConditions = append(andConditions, res)
					}
				}
			}
			if len(andConditions) > 0 {
				q = fmt.Sprintf("%s where %s", q, strings.Join(andConditions, " and "))
			}

			//sort
			sortDef := func(sort srequest.SortItem) string {
				switch sort.Field {
				case "first_name":
					fallthrough
				case "last_name":
					{
						return fmt.Sprintf("%s %s", sort.Field, sort.Dir)
					}
				}
				return ""
			}

			return q, sortDef, params, nil
		}

	}
	if len(andConditions) > 0 {
		q = fmt.Sprintf("%s where %s", q, strings.Join(andConditions, " and "))
	}

	//sort
	sortDef := func(sort srequest.SortItem) string {
		switch sort.Field {
		case "first_name":
			fallthrough
		case "last_name":
			{
				return fmt.Sprintf("%s %s", sort.Field, sort.Dir)
			}
		}
		return ""
	}

	return q, sortDef, params, nil
}
