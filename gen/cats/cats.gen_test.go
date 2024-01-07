package cats

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/testcontainers/testcontainers-go"
	tc "github.com/testcontainers/testcontainers-go/modules/nats"
)

type MockHandler struct{}

func (m MockHandler) NewCat(r micro.Request, cat CatIntake)                {}
func (m MockHandler) EditCat(r micro.Request, catID string, cat CatIntake) {}
func (m MockHandler) GetCat(r micro.Request, catID string)                 {}

func TestCatService(t *testing.T) {
	ctx := context.Background()

	natsContainer, err := tc.RunContainer(ctx,
		testcontainers.WithImage("nats:2"),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Clean up the container
	defer func() {
		if err := natsContainer.Terminate(ctx); err != nil {
			t.Fatalf(err.Error())
		}
	}()

	connectionURL, err := natsContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	nc, err := nats.Connect(connectionURL)
	if err != nil {
		t.Fatalf(err.Error())
	}

	h := MockHandler{}
	if _, err := CreateService(nc, h, &Options{
		Name:    "CatsService",
		Version: "0.0.1",
	}); err != nil {
		t.Fatalf("err creating service: %s\n", err)
	}
}
