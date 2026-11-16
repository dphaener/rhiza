package index

import "github.com/dphaener/rhiza/internal/model"

// AsIs returns the as-is view, as model.AsIs computes it over the change
// sets folded from every stored entry, assembled from the stored as-is
// rows: every node with its fields and history, every link, and the as-is
// problems. It decodes no log entry. The duplicate-handle problems are
// the model's own, of the view's final state, as model.AsIs reports them.
// The error is the store's: a closed index returns one.
func (x *Index) AsIs() (model.Model, []model.Problem, error) {
	tx, err := x.db.Begin()
	if err != nil {
		return model.Model{}, nil, err
	}
	defer tx.Rollback()

	seed, err := loadSeed(tx)
	if err != nil {
		return model.Model{}, nil, err
	}
	// To-be over the seed with no accepted change set layered is the seed
	// itself, finished as model.AsIs finishes its view.
	m, ps := model.ToBeFrom(seed, nil)
	return m, ps, nil
}
