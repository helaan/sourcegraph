package graphqlbackend

import (
	"context"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/backend"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend/graphqlutil"
	"github.com/sourcegraph/sourcegraph/pkg/api"
	"github.com/sourcegraph/sourcegraph/pkg/vcs/git"
)

func (r *gitTreeEntryResolver) IsRoot() bool {
	path := path.Clean(r.path)
	return path == "/" || path == "." || path == ""
}

type gitTreeEntryConnectionArgs struct {
	graphqlutil.ConnectionArgs
	Recursive bool
	// If recurseSingleChild is true, we will return a flat list of every
	// directory and file in a single-child nest.
	RecursiveSingleChild bool
}

func (r *gitTreeEntryResolver) Entries(ctx context.Context, args *gitTreeEntryConnectionArgs) ([]*gitTreeEntryResolver, error) {
	return r.entries(ctx, args, nil)
}

func (r *gitTreeEntryResolver) Directories(ctx context.Context, args *gitTreeEntryConnectionArgs) ([]*gitTreeEntryResolver, error) {
	return r.entries(ctx, args, func(fi os.FileInfo) bool { return fi.Mode().IsDir() })
}

func (r *gitTreeEntryResolver) Files(ctx context.Context, args *gitTreeEntryConnectionArgs) ([]*gitTreeEntryResolver, error) {
	return r.entries(ctx, args, func(fi os.FileInfo) bool { return !fi.Mode().IsDir() })
}

func (r *gitTreeEntryResolver) entries(ctx context.Context, args *gitTreeEntryConnectionArgs, filter func(fi os.FileInfo) bool) ([]*gitTreeEntryResolver, error) {
	cachedRepo, err := backend.CachedGitRepo(ctx, r.commit.repo.repo)
	if err != nil {
		return nil, err
	}
	entries, err := git.ReadDir(ctx, *cachedRepo, api.CommitID(r.commit.OID()), r.path, r.isRecursive || args.Recursive)
	if err != nil {
		if strings.Contains(err.Error(), "file does not exist") { // TODO proper error value
			// empty tree is not an error
		} else {
			return nil, err
		}
	}

	sort.Sort(byDirectory(entries))

	if args.First != nil && len(entries) > int(*args.First) {
		entries = entries[:int(*args.First)]
	}

	var prefix string
	if r.path != "" {
		prefix = r.path + "/"
	}

	var l []*gitTreeEntryResolver
	for _, entry := range entries {
		if filter == nil || filter(entry) {
			l = append(l, &gitTreeEntryResolver{
				commit: r.commit,
				path:   prefix + entry.Name(), // relies on git paths being cleaned already
				stat:   entry,
			})
		}
	}

	if !args.Recursive && args.RecursiveSingleChild && len(l) == 1 {
		subEntries, err := l[0].entries(ctx, args, filter)
		if err != nil {
			return nil, err
		}
		l = append(l, subEntries...)
	}

	return l, nil
}

type byDirectory []os.FileInfo

func (s byDirectory) Len() int {
	return len(s)
}

func (s byDirectory) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byDirectory) Less(i, j int) bool {
	if s[i].IsDir() && !s[j].IsDir() {
		return true
	}

	if !s[i].IsDir() && s[j].IsDir() {
		return false
	}

	return s[i].Name() < s[j].Name()
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_147(size int) error {
	const bufSize = 1024

	f, err := os.Create("/tmp/test")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fb := bufio.NewWriter(f)
	defer fb.Flush()

	buf := make([]byte, bufSize)

	for i := size; i > 0; i -= bufSize {
		if _, err = rand.Read(buf); err != nil {
			fmt.Printf("error occurred during random: %!s(MISSING)\n", err)
			break
		}
		bR := bytes.NewReader(buf)
		if _, err = io.Copy(fb, bR); err != nil {
			fmt.Printf("failed during copy: %!s(MISSING)\n", err)
			break
		}
	}

	return err
}		
