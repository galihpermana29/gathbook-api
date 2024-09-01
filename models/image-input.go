package models

import "mime/multipart"

type UploadImageInput struct {
	Images []*multipart.FileHeader `form:"images" binding:"required"`
}
