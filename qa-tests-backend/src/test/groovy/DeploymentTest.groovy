import groups.BAT
import io.stackrox.proto.api.v1.SearchServiceOuterClass.RawQuery
import org.junit.experimental.categories.Category
import services.DeploymentService
import services.ImageService
import objects.Deployment
import spock.lang.Unroll
import util.Timer

class DeploymentTest extends BaseSpecification {
    private static final String DEPLOYMENT_NAME = "image-join"
    private static final Deployment DEPLOYMENT = new Deployment()
            .setName(DEPLOYMENT_NAME)
            .setImage("nginx@sha256:204a9a8e65061b10b92ad361dd6f406248404fe60efd5d6a8f2595f18bb37aad")
            .addLabel("app", "test")
            .setCommand(["sh", "-c", "apt-get -y update && sleep 600"])

    def setupSpec() {
        orchestrator.createDeployment(DEPLOYMENT)
    }

    def cleanupSpec() {
        orchestrator.deleteDeployment(DEPLOYMENT)
    }

    @Unroll
    @Category([BAT])
    def "Verify deployment -> image links #query"() {
        when:
        Timer t = new Timer(3, 10)
        def img = null
        while (img == null && t.IsValid()) {
            img = ImageService.getImage(
                    "sha256:204a9a8e65061b10b92ad361dd6f406248404fe60efd5d6a8f2595f18bb37aad", false)
        }
        assert img != null

        then:
        def results = DeploymentService.listDeploymentsSearch(RawQuery.newBuilder().setQuery(query).build())
        assert results.deploymentsList.find { x -> x.getName() == DEPLOYMENT_NAME } != null

        where:
        "Data inputs are: "
        query                                                                                                   | _
        "Image:docker.io/library/nginx@sha256:204a9a8e65061b10b92ad361dd6f406248404fe60efd5d6a8f2595f18bb37aad" | _
        "Image Sha:sha256:204a9a8e65061b10b92ad361dd6f406248404fe60efd5d6a8f2595f18bb37aad"                     | _
        "CVE:CVE-2018-18314+Fixable:true"                                                                       | _
        "Deployment:${DEPLOYMENT_NAME}+Image:r/docker.*"                                                        | _
        "Image:r/docker.*"                                                                                      | _
        "Image:!stackrox.io"                                                                                    | _
        "Deployment:${DEPLOYMENT_NAME}+Image:!stackrox.io"                                                      | _
        "Image Remote:library/nginx+Image Registry:docker.io"                                                   | _
    }

    @Unroll
    @Category([BAT])
    def "Verify image -> deployment links #query"() {
        when:
        Timer t = new Timer(3, 10)
        def img = null
        while (img == null && t.IsValid()) {
            img = ImageService.getImage(
                    "sha256:204a9a8e65061b10b92ad361dd6f406248404fe60efd5d6a8f2595f18bb37aad", false)
        }
        assert img != null

        then:
        def images = ImageService.getImages(RawQuery.newBuilder().setQuery(query).build())
        assert images.find {
            x -> x.getId() == "sha256:204a9a8e65061b10b92ad361dd6f406248404fe60efd5d6a8f2595f18bb37aad" } != null

        where:
        "Data inputs are: "
        query                                                                                                   | _
        "Deployment:${DEPLOYMENT_NAME}"                                                                         | _
        "Label:app=test"                                                                                        | _
        "Image:docker.io/library/nginx@sha256:204a9a8e65061b10b92ad361dd6f406248404fe60efd5d6a8f2595f18bb37aad" | _
        "Label:app=test+Image:docker.io/library/nginx"                                                          | _
    }

}
