package converter

import (
	"context"
	"github.com/castaneai/sankaku"
	"net/http"
	"testing"
)

func TestSankakuConverter_Convert(t *testing.T) {
	hc := &http.Client{}
	s, err := sankaku.NewClient(hc, "https://chan.sankakucomplex.com", "en", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	sp, sd, err := s.GetPostWithDetail(ctx, "12345")
	if err != nil {
		t.Fatal(err)
	}

	cv, err := NewSankakuConverter()
	if err != nil {
		t.Fatal(err)
	}

	post, err := cv.Convert(sp, sd)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", post)
}
