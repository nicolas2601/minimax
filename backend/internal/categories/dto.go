package categories

type CreateRequest struct {
	Name     string  `json:"name" binding:"required,min=1,max=100"`
	Type     string  `json:"type" binding:"required,oneof=expense income"`
	ParentID *string `json:"parent_id,omitempty"`
	Icon     *string `json:"icon,omitempty" binding:"omitempty,max=50"`
	Color    *string `json:"color,omitempty" binding:"omitempty,len=7"`
}

type UpdateRequest struct {
	Name  *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Color *string `json:"color,omitempty" binding:"omitempty,len=7"`
	Icon  *string `json:"icon,omitempty" binding:"omitempty,max=50"`
}

type ListResponse struct {
	Categories []Category `json:"categories"`
}