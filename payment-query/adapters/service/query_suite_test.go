// +build test

package service_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	. "github.com/fractalpal/eventflow-example/payment-query/adapters/service"
	"github.com/fractalpal/eventflow-example/payment-query/app"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
	Repository                    *MockedRepository
	Service                       app.Service
}

func (suite *ServiceTestSuite) SetupTest() {
	//suite.VariableThatShouldStartAtFive = 5
	suite.Repository = new(MockedRepository)
	suite.Service = NewQuery(suite.Repository)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) Test_FindByIDSuccess() {
	var payment app.Payment
	id := "ee3e97ba-04bf-4422-a112-08c6ca89effb"
	ctx := context.TODO()
	payment.ID = id
	suite.Repository.On("FindByID", ctx, id).Return(payment, nil)

	created, err := suite.Service.FindByID(ctx, id)

	assert.NoError(suite.T(), err, "expected no error in service find by ID")
	assert.Equal(suite.T(), id, created.ID)
}

func (suite *ServiceTestSuite) Test_FindAllSuccess() {
	var payment app.Payment
	ctx := context.TODO()
	suite.Repository.On("FindAll", ctx).Return([]app.Payment{payment}, nil)

	result, err := suite.Service.FindAll(ctx, int64(1), int64(0))

	assert.NoError(suite.T(), err, "expected no error in service find all")
	assert.Equal(suite.T(), 1, len(result))
}
