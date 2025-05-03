package utils

import (
	"net/http"

	"github.com/srgjo27/e-learning/internal/entity"
)

func RBACMiddleware(allowedRoles ...entity.Role) func(http.Handler) http.Handler {
	roleSet := make(map[entity.Role]struct{})
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value("userRole")
			if roleVal == nil {
				http.Error(w, "Forbidden - role missing", http.StatusForbidden)
				return
			}

			role, ok := roleVal.(entity.Role)
			if !ok {
				if strRole, ok := roleVal.(string); ok {
					role = entity.Role(strRole)
				} else {
					http.Error(w, "Forbidden - invalid role", http.StatusForbidden)
					return
				}
			}

			if _, allowed := roleSet[role]; !allowed {
				http.Error(w, "Forbidden - insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}