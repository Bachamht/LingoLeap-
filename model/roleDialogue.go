
package model
import (

)

type RoleDialaogue_req struct {
	Role string `json:"role"`
	Message string `json:"message"`
	User_id     int `json:"user_id"`
}

type RoleDialaogue_res struct {
	Message string `json:"message"`
	Session_id  int

}

func CreateNewRoleDialaogue(requestData *RoleDialaogue_req) (string, error){
	roleDialaogue := requestData.Role
	userID := requestData.User_id
	prompt := "你现在扮演的角色是" + roleDialaogue
	role := "system"
	answer, err1 := IteracionWithAI(prompt, userID, role)
    if err1 != nil {
        return "", nil
    }
	return answer, nil
}