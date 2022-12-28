package clouduploader

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/otterize/network-mapper/src/mapper/pkg/cloudclient"
	cloudclientmocks "github.com/otterize/network-mapper/src/mapper/pkg/cloudclient/mocks"
	"github.com/otterize/network-mapper/src/mapper/pkg/config"
	"github.com/otterize/network-mapper/src/mapper/pkg/graph/model"
	"github.com/otterize/network-mapper/src/mapper/pkg/resolvers"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
	"testing"
)

type CloudUploaderTestSuite struct {
	suite.Suite
	testNamespace string
	intentsHolder *resolvers.IntentsHolder
	cloudUploader *CloudUploader
	cloudConfig   Config
	clientMock    *cloudclientmocks.MockCloudClient
}

func (s *CloudUploaderTestSuite) SetupTest() {
	s.testNamespace = "test-namespace"
	s.intentsHolder = resolvers.NewIntentsHolder(nil, resolvers.IntentsHolderConfig{StoreConfigMap: config.StoreConfigMapDefault, Namespace: s.testNamespace})
	s.cloudConfig = Config{
		ClientId: "test-client-id",
	}
}

func (s *CloudUploaderTestSuite) BeforeTest(_, testName string) {
	controller := gomock.NewController(s.T())
	factory := s.GetCloudClientFactoryMock(controller)
	s.cloudUploader = NewCloudUploader(s.intentsHolder, s.cloudConfig, factory)
}

func (s *CloudUploaderTestSuite) GetCloudClientFactoryMock(controller *gomock.Controller) cloudclient.FactoryFunction {
	s.clientMock = cloudclientmocks.NewMockCloudClient(controller)

	factory := func(ctx context.Context, apiAddress string, tokenSource oauth2.TokenSource) cloudclient.CloudClient {
		return s.clientMock
	}
	return factory
}

func (s *CloudUploaderTestSuite) addIntent(source string, destination string) {
	s.intentsHolder.AddIntent(
		model.OtterizeServiceIdentity{Name: source, Namespace: s.testNamespace},
		model.OtterizeServiceIdentity{Name: destination, Namespace: s.testNamespace},
	)
}

func (s *CloudUploaderTestSuite) TestUploadIntents() {
	s.addIntent("client1", "server1")
	s.addIntent("client1", "server2")

	intents1 := []cloudclient.IntentInput{
		{ClientName: "client1", ServerName: "server1", Namespace: s.testNamespace},
		{ClientName: "client1", ServerName: "server2", Namespace: s.testNamespace},
	}
	s.clientMock.EXPECT().ReportDiscoveredIntents(gomock.InAnyOrder(intents1)).Return(true).Times(1)

	s.cloudUploader.uploadDiscoveredIntents(context.Background())

	s.addIntent("client2", "server1")

	intents2 := []cloudclient.IntentInput{
		{ClientName: "client1", ServerName: "server1", Namespace: s.testNamespace},
		{ClientName: "client1", ServerName: "server2", Namespace: s.testNamespace},
		{ClientName: "client2", ServerName: "server1", Namespace: s.testNamespace},
	}

	s.clientMock.EXPECT().ReportDiscoveredIntents(gomock.InAnyOrder(intents2)).Return(true).Times(1)
	s.cloudUploader.uploadDiscoveredIntents(context.Background())
}

func (s *CloudUploaderTestSuite) TestDontUploadWithoutIntents() {
	s.clientMock.EXPECT().ReportDiscoveredIntents(gomock.Any()).Times(0)

	s.cloudUploader.uploadDiscoveredIntents(context.Background())
}

func (s *CloudUploaderTestSuite) TestUploadSameIntentOnce() {
	s.addIntent("client", "server")

	intents := []cloudclient.IntentInput{
		{ClientName: "client", ServerName: "server", Namespace: s.testNamespace},
	}

	s.clientMock.EXPECT().ReportDiscoveredIntents(intents).Return(true).Times(1)

	s.cloudUploader.uploadDiscoveredIntents(context.Background())
	s.addIntent("client", "server")
	s.cloudUploader.uploadDiscoveredIntents(context.Background())
}

func (s *CloudUploaderTestSuite) TestRetryOnFailed() {
	s.addIntent("client", "server")

	intents := []cloudclient.IntentInput{
		{ClientName: "client", ServerName: "server", Namespace: s.testNamespace},
	}

	s.clientMock.EXPECT().ReportDiscoveredIntents(intents).Return(false).Times(1)

	s.clientMock.EXPECT().ReportDiscoveredIntents(intents).Return(true).Times(1)

	s.cloudUploader.uploadDiscoveredIntents(context.Background())
	s.cloudUploader.uploadDiscoveredIntents(context.Background())
}

func (s *CloudUploaderTestSuite) TestDontUploadWhenNothingNew() {
	s.addIntent("client", "server")

	intents := []cloudclient.IntentInput{
		{ClientName: "client", ServerName: "server", Namespace: s.testNamespace},
	}

	s.clientMock.EXPECT().ReportDiscoveredIntents(intents).Return(true).Times(1)

	s.cloudUploader.uploadDiscoveredIntents(context.Background())
	s.cloudUploader.uploadDiscoveredIntents(context.Background())
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(CloudUploaderTestSuite))
}
