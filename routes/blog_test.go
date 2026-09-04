package routes

import (
	"testing"

	"fiber-ecommerce/models"
)

func TestValidateBlogRequiresEveryLanguage(t *testing.T) {
	tests := []struct {
		name    string
		blog    models.Blog
		wantErr bool
	}{
		{
			name: "complete blog",
			blog: models.Blog{
				Image: "/uploads/cover.webp",
				Title: "English title***Русский заголовок***O‘zbekcha sarlavha",
				Text:  "<p>English text</p>***<p>Русский текст</p>***<p>O‘zbekcha matn</p>",
			},
		},
		{
			name: "missing localized title",
			blog: models.Blog{
				Image: "/uploads/cover.webp",
				Title: "English title******O‘zbekcha sarlavha",
				Text:  "<p>English text</p>***<p>Русский текст</p>***<p>O‘zbekcha matn</p>",
			},
			wantErr: true,
		},
		{
			name: "empty rich text markup",
			blog: models.Blog{
				Image: "/uploads/cover.webp",
				Title: "English title***Русский заголовок***O‘zbekcha sarlavha",
				Text:  "<p><br></p>***<p>Русский текст</p>***<p>O‘zbekcha matn</p>",
			},
			wantErr: true,
		},
		{
			name: "missing image",
			blog: models.Blog{
				Title: "English title***Русский заголовок***O‘zbekcha sarlavha",
				Text:  "<p>English text</p>***<p>Русский текст</p>***<p>O‘zbekcha matn</p>",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateBlog(test.blog)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateBlog() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
