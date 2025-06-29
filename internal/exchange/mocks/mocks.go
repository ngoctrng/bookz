// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"github.com/google/uuid"
	"github.com/ngoctrng/bookz/internal/exchange"
	"github.com/ngoctrng/bookz/internal/exchange/usecases"
	mock "github.com/stretchr/testify/mock"
)

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

type MockRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRepository) EXPECT() *MockRepository_Expecter {
	return &MockRepository_Expecter{mock: &_m.Mock}
}

// FetchRequestedBookOwner provides a mock function for the type MockRepository
func (_mock *MockRepository) FetchRequestedBookOwner(id int) (uuid.UUID, error) {
	ret := _mock.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for FetchRequestedBookOwner")
	}

	var r0 uuid.UUID
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(int) (uuid.UUID, error)); ok {
		return returnFunc(id)
	}
	if returnFunc, ok := ret.Get(0).(func(int) uuid.UUID); ok {
		r0 = returnFunc(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(int) error); ok {
		r1 = returnFunc(id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_FetchRequestedBookOwner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchRequestedBookOwner'
type MockRepository_FetchRequestedBookOwner_Call struct {
	*mock.Call
}

// FetchRequestedBookOwner is a helper method to define mock.On call
//   - id int
func (_e *MockRepository_Expecter) FetchRequestedBookOwner(id interface{}) *MockRepository_FetchRequestedBookOwner_Call {
	return &MockRepository_FetchRequestedBookOwner_Call{Call: _e.mock.On("FetchRequestedBookOwner", id)}
}

func (_c *MockRepository_FetchRequestedBookOwner_Call) Run(run func(id int)) *MockRepository_FetchRequestedBookOwner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 int
		if args[0] != nil {
			arg0 = args[0].(int)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRepository_FetchRequestedBookOwner_Call) Return(uUID uuid.UUID, err error) *MockRepository_FetchRequestedBookOwner_Call {
	_c.Call.Return(uUID, err)
	return _c
}

func (_c *MockRepository_FetchRequestedBookOwner_Call) RunAndReturn(run func(id int) (uuid.UUID, error)) *MockRepository_FetchRequestedBookOwner_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function for the type MockRepository
func (_mock *MockRepository) GetAll(uid uuid.UUID) ([]*exchange.Proposal, error) {
	ret := _mock.Called(uid)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*exchange.Proposal
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) ([]*exchange.Proposal, error)); ok {
		return returnFunc(uid)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) []*exchange.Proposal); ok {
		r0 = returnFunc(uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*exchange.Proposal)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = returnFunc(uid)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - uid uuid.UUID
func (_e *MockRepository_Expecter) GetAll(uid interface{}) *MockRepository_GetAll_Call {
	return &MockRepository_GetAll_Call{Call: _e.mock.On("GetAll", uid)}
}

func (_c *MockRepository_GetAll_Call) Run(run func(uid uuid.UUID)) *MockRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 uuid.UUID
		if args[0] != nil {
			arg0 = args[0].(uuid.UUID)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRepository_GetAll_Call) Return(proposals []*exchange.Proposal, err error) *MockRepository_GetAll_Call {
	_c.Call.Return(proposals, err)
	return _c
}

func (_c *MockRepository_GetAll_Call) RunAndReturn(run func(uid uuid.UUID) ([]*exchange.Proposal, error)) *MockRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function for the type MockRepository
func (_mock *MockRepository) GetByID(id int) (*exchange.Proposal, error) {
	ret := _mock.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *exchange.Proposal
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(int) (*exchange.Proposal, error)); ok {
		return returnFunc(id)
	}
	if returnFunc, ok := ret.Get(0).(func(int) *exchange.Proposal); ok {
		r0 = returnFunc(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*exchange.Proposal)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(int) error); ok {
		r1 = returnFunc(id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - id int
func (_e *MockRepository_Expecter) GetByID(id interface{}) *MockRepository_GetByID_Call {
	return &MockRepository_GetByID_Call{Call: _e.mock.On("GetByID", id)}
}

func (_c *MockRepository_GetByID_Call) Run(run func(id int)) *MockRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 int
		if args[0] != nil {
			arg0 = args[0].(int)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRepository_GetByID_Call) Return(proposal *exchange.Proposal, err error) *MockRepository_GetByID_Call {
	_c.Call.Return(proposal, err)
	return _c
}

func (_c *MockRepository_GetByID_Call) RunAndReturn(run func(id int) (*exchange.Proposal, error)) *MockRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function for the type MockRepository
func (_mock *MockRepository) Save(p *exchange.Proposal) error {
	ret := _mock.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(*exchange.Proposal) error); ok {
		r0 = returnFunc(p)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - p *exchange.Proposal
func (_e *MockRepository_Expecter) Save(p interface{}) *MockRepository_Save_Call {
	return &MockRepository_Save_Call{Call: _e.mock.On("Save", p)}
}

func (_c *MockRepository_Save_Call) Run(run func(p *exchange.Proposal)) *MockRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 *exchange.Proposal
		if args[0] != nil {
			arg0 = args[0].(*exchange.Proposal)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockRepository_Save_Call) Return(err error) *MockRepository_Save_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRepository_Save_Call) RunAndReturn(run func(p *exchange.Proposal) error) *MockRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMessageBus creates a new instance of MockMessageBus. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMessageBus(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMessageBus {
	mock := &MockMessageBus{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockMessageBus is an autogenerated mock type for the MessageBus type
type MockMessageBus struct {
	mock.Mock
}

type MockMessageBus_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMessageBus) EXPECT() *MockMessageBus_Expecter {
	return &MockMessageBus_Expecter{mock: &_m.Mock}
}

// PublishProposalAccepted provides a mock function for the type MockMessageBus
func (_mock *MockMessageBus) PublishProposalAccepted(p *exchange.Proposal) error {
	ret := _mock.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for PublishProposalAccepted")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(*exchange.Proposal) error); ok {
		r0 = returnFunc(p)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockMessageBus_PublishProposalAccepted_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishProposalAccepted'
type MockMessageBus_PublishProposalAccepted_Call struct {
	*mock.Call
}

// PublishProposalAccepted is a helper method to define mock.On call
//   - p *exchange.Proposal
func (_e *MockMessageBus_Expecter) PublishProposalAccepted(p interface{}) *MockMessageBus_PublishProposalAccepted_Call {
	return &MockMessageBus_PublishProposalAccepted_Call{Call: _e.mock.On("PublishProposalAccepted", p)}
}

func (_c *MockMessageBus_PublishProposalAccepted_Call) Run(run func(p *exchange.Proposal)) *MockMessageBus_PublishProposalAccepted_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 *exchange.Proposal
		if args[0] != nil {
			arg0 = args[0].(*exchange.Proposal)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockMessageBus_PublishProposalAccepted_Call) Return(err error) *MockMessageBus_PublishProposalAccepted_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockMessageBus_PublishProposalAccepted_Call) RunAndReturn(run func(p *exchange.Proposal) error) *MockMessageBus_PublishProposalAccepted_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUsecase creates a new instance of MockUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUsecase {
	mock := &MockUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockUsecase is an autogenerated mock type for the Usecase type
type MockUsecase struct {
	mock.Mock
}

type MockUsecase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUsecase) EXPECT() *MockUsecase_Expecter {
	return &MockUsecase_Expecter{mock: &_m.Mock}
}

// AcceptProposal provides a mock function for the type MockUsecase
func (_mock *MockUsecase) AcceptProposal(id int, uid uuid.UUID) error {
	ret := _mock.Called(id, uid)

	if len(ret) == 0 {
		panic("no return value specified for AcceptProposal")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(int, uuid.UUID) error); ok {
		r0 = returnFunc(id, uid)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockUsecase_AcceptProposal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AcceptProposal'
type MockUsecase_AcceptProposal_Call struct {
	*mock.Call
}

// AcceptProposal is a helper method to define mock.On call
//   - id int
//   - uid uuid.UUID
func (_e *MockUsecase_Expecter) AcceptProposal(id interface{}, uid interface{}) *MockUsecase_AcceptProposal_Call {
	return &MockUsecase_AcceptProposal_Call{Call: _e.mock.On("AcceptProposal", id, uid)}
}

func (_c *MockUsecase_AcceptProposal_Call) Run(run func(id int, uid uuid.UUID)) *MockUsecase_AcceptProposal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 int
		if args[0] != nil {
			arg0 = args[0].(int)
		}
		var arg1 uuid.UUID
		if args[1] != nil {
			arg1 = args[1].(uuid.UUID)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockUsecase_AcceptProposal_Call) Return(err error) *MockUsecase_AcceptProposal_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockUsecase_AcceptProposal_Call) RunAndReturn(run func(id int, uid uuid.UUID) error) *MockUsecase_AcceptProposal_Call {
	_c.Call.Return(run)
	return _c
}

// CreateProposal provides a mock function for the type MockUsecase
func (_mock *MockUsecase) CreateProposal(in usecases.CreateProposalInput) error {
	ret := _mock.Called(in)

	if len(ret) == 0 {
		panic("no return value specified for CreateProposal")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(usecases.CreateProposalInput) error); ok {
		r0 = returnFunc(in)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockUsecase_CreateProposal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateProposal'
type MockUsecase_CreateProposal_Call struct {
	*mock.Call
}

// CreateProposal is a helper method to define mock.On call
//   - in usecases.CreateProposalInput
func (_e *MockUsecase_Expecter) CreateProposal(in interface{}) *MockUsecase_CreateProposal_Call {
	return &MockUsecase_CreateProposal_Call{Call: _e.mock.On("CreateProposal", in)}
}

func (_c *MockUsecase_CreateProposal_Call) Run(run func(in usecases.CreateProposalInput)) *MockUsecase_CreateProposal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 usecases.CreateProposalInput
		if args[0] != nil {
			arg0 = args[0].(usecases.CreateProposalInput)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockUsecase_CreateProposal_Call) Return(err error) *MockUsecase_CreateProposal_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockUsecase_CreateProposal_Call) RunAndReturn(run func(in usecases.CreateProposalInput) error) *MockUsecase_CreateProposal_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllProposals provides a mock function for the type MockUsecase
func (_mock *MockUsecase) GetAllProposals(uid uuid.UUID) ([]*exchange.Proposal, error) {
	ret := _mock.Called(uid)

	if len(ret) == 0 {
		panic("no return value specified for GetAllProposals")
	}

	var r0 []*exchange.Proposal
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) ([]*exchange.Proposal, error)); ok {
		return returnFunc(uid)
	}
	if returnFunc, ok := ret.Get(0).(func(uuid.UUID) []*exchange.Proposal); ok {
		r0 = returnFunc(uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*exchange.Proposal)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = returnFunc(uid)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockUsecase_GetAllProposals_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllProposals'
type MockUsecase_GetAllProposals_Call struct {
	*mock.Call
}

// GetAllProposals is a helper method to define mock.On call
//   - uid uuid.UUID
func (_e *MockUsecase_Expecter) GetAllProposals(uid interface{}) *MockUsecase_GetAllProposals_Call {
	return &MockUsecase_GetAllProposals_Call{Call: _e.mock.On("GetAllProposals", uid)}
}

func (_c *MockUsecase_GetAllProposals_Call) Run(run func(uid uuid.UUID)) *MockUsecase_GetAllProposals_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 uuid.UUID
		if args[0] != nil {
			arg0 = args[0].(uuid.UUID)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockUsecase_GetAllProposals_Call) Return(proposals []*exchange.Proposal, err error) *MockUsecase_GetAllProposals_Call {
	_c.Call.Return(proposals, err)
	return _c
}

func (_c *MockUsecase_GetAllProposals_Call) RunAndReturn(run func(uid uuid.UUID) ([]*exchange.Proposal, error)) *MockUsecase_GetAllProposals_Call {
	_c.Call.Return(run)
	return _c
}

// GetProposalByID provides a mock function for the type MockUsecase
func (_mock *MockUsecase) GetProposalByID(id int) (*exchange.Proposal, error) {
	ret := _mock.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetProposalByID")
	}

	var r0 *exchange.Proposal
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(int) (*exchange.Proposal, error)); ok {
		return returnFunc(id)
	}
	if returnFunc, ok := ret.Get(0).(func(int) *exchange.Proposal); ok {
		r0 = returnFunc(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*exchange.Proposal)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(int) error); ok {
		r1 = returnFunc(id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockUsecase_GetProposalByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProposalByID'
type MockUsecase_GetProposalByID_Call struct {
	*mock.Call
}

// GetProposalByID is a helper method to define mock.On call
//   - id int
func (_e *MockUsecase_Expecter) GetProposalByID(id interface{}) *MockUsecase_GetProposalByID_Call {
	return &MockUsecase_GetProposalByID_Call{Call: _e.mock.On("GetProposalByID", id)}
}

func (_c *MockUsecase_GetProposalByID_Call) Run(run func(id int)) *MockUsecase_GetProposalByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 int
		if args[0] != nil {
			arg0 = args[0].(int)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockUsecase_GetProposalByID_Call) Return(proposal *exchange.Proposal, err error) *MockUsecase_GetProposalByID_Call {
	_c.Call.Return(proposal, err)
	return _c
}

func (_c *MockUsecase_GetProposalByID_Call) RunAndReturn(run func(id int) (*exchange.Proposal, error)) *MockUsecase_GetProposalByID_Call {
	_c.Call.Return(run)
	return _c
}
