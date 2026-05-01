package content

import "io/fs"

type PostRepository interface {
	All() ([]Post, error)
	Published() ([]Post, error)
	BySlug(slug string) (Post, error)
	ByTag(tag string) ([]Post, error)
}

type FilePostRepository struct {
	fileSys fs.FS
	cache   map[string]Post
}

func NewFilePostRepository(fileSys fs.FS) (*FilePostRepository, error) {
	r := &FilePostRepository{
		fileSys: fileSys,
		cache:   map[string]Post{},
	}

	if err := r.Load(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *FilePostRepository) All() ([]Post, error) {

}
