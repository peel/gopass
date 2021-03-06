package root

import (
	"context"
	"sort"
	"testing"

	"path"

	"github.com/justwatchcom/gopass/backend"
	"github.com/justwatchcom/gopass/config"
	"github.com/justwatchcom/gopass/tests/gptest"
	"github.com/stretchr/testify/assert"
)

func TestSimpleList(t *testing.T) {
	ctx := context.Background()

	u := gptest.NewUnitTester(t)
	defer u.Remove()

	rs, err := createRootStore(ctx, u)
	assert.NoError(t, err)

	tree, err := rs.Tree(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo"}, tree.List(0))
}

func TestListMulti(t *testing.T) {
	ctx := context.Background()

	u := gptest.NewUnitTester(t)
	defer u.Remove()

	// root store
	rs, err := createRootStore(ctx, u)
	assert.NoError(t, err)

	ents := make([]string, 0, 3*len(u.Entries))
	ents = append(ents, u.Entries...)

	// sub1 store
	assert.NoError(t, u.InitStore("sub1"))
	for _, k := range u.Entries {
		ents = append(ents, path.Join("sub1", k))
	}

	// sub2 store
	assert.NoError(t, u.InitStore("sub2"))
	for _, k := range u.Entries {
		ents = append(ents, path.Join("sub2", k))
	}

	assert.NoError(t, rs.AddMount(ctx, "sub1", u.StoreDir("sub1")))
	assert.NoError(t, rs.AddMount(ctx, "sub2", u.StoreDir("sub2")))

	tree, err := rs.Tree(ctx)
	assert.NoError(t, err)

	sort.Strings(ents)
	lst := tree.List(0)
	sort.Strings(lst)
	assert.Equal(t, ents, lst)
}

func TestListNested(t *testing.T) {
	ctx := context.Background()

	u := gptest.NewUnitTester(t)
	defer u.Remove()

	// root store
	rs, err := createRootStore(ctx, u)
	assert.NoError(t, err)

	ents := make([]string, 0, 3*len(u.Entries))
	ents = append(ents, u.Entries...)

	// sub1 store
	assert.NoError(t, u.InitStore("sub1"))
	for _, k := range u.Entries {
		ents = append(ents, path.Join("sub1", k))
	}

	// sub2 store
	assert.NoError(t, u.InitStore("sub2"))
	for _, k := range u.Entries {
		ents = append(ents, path.Join("sub2", k))
	}

	// sub3 store
	assert.NoError(t, u.InitStore("sub3"))
	for _, k := range u.Entries {
		ents = append(ents, path.Join("sub2", "sub3", k))
	}

	assert.NoError(t, rs.AddMount(ctx, "sub1", u.StoreDir("sub1")))
	assert.NoError(t, rs.AddMount(ctx, "sub2", u.StoreDir("sub2")))
	assert.NoError(t, rs.AddMount(ctx, "sub2/sub3", u.StoreDir("sub3")))

	tree, err := rs.Tree(ctx)
	assert.NoError(t, err)

	sort.Strings(ents)
	lst := tree.List(0)
	sort.Strings(lst)
	assert.Equal(t, ents, lst)

	assert.Equal(t, false, rs.Exists(ctx, "sub1"))
	// TODO create entry and text Exists == true
	assert.Equal(t, true, rs.IsDir(ctx, "sub1"))
	//assert.Equal(t, "", rs.String()) // TODO
	assert.Equal(t, "", rs.Alias())
	assert.NotNil(t, rs.Store(ctx, "sub1"))
}

func createRootStore(ctx context.Context, u *gptest.Unit) (*Store, error) {
	ctx = backend.WithSyncBackendString(ctx, "gitmock")
	ctx = backend.WithCryptoBackendString(ctx, "gpgmock")
	return New(
		ctx,
		&config.Config{
			Root: &config.StoreConfig{
				Path: u.StoreDir(""),
			},
		},
	)
}
