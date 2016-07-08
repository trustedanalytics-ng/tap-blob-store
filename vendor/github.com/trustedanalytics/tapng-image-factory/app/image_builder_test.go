package main

import (
	"testing"

	"archive/tar"
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
)

var (
	testBaseImage            = ImagesMap["JAVA"]
	testImageArtifactContent = "Test artifact file content"
	testDockerfileContent    = "FROM " + testBaseImage + "\n" + "RUN mkdir /test_dir\n" + "CMD [~/run.sh]"

	testImageArtifact = bytes.NewBufferString(testImageArtifactContent)
	testDockerfile    = bytes.NewBufferString(testDockerfileContent)
)

func TestCreateDockerfile(t *testing.T) {
	Convey("Test CreateDockerfile", t, func() {
		Convey("Dockerfile should contains proper lines", func() {
			dockerfile, _ := createDockerfile("JAVA")
			dockerfileStr, _ := StreamToString(dockerfile)
			dockerfileStringArray := strings.Split(dockerfileStr, "\n")

			So(dockerfileStringArray[0], ShouldEqual, "FROM "+GetHubAddress()+"/"+testBaseImage)
			So(dockerfileStringArray[1], ShouldEqual, "ADD "+artifactFileName+" /root")
			So(dockerfileStringArray[2], ShouldEqual, "CMD [/root/run.sh]")
		})
	})
}

func TestCreateBuildContext(t *testing.T) {
	Convey("Test CreateBuildContext", t, func() {
		Convey("Context should contains proper files", func() {
			context, err := createBuildContext(testImageArtifact, testDockerfile)
			So(err, ShouldBeNil)
			tr := tar.NewReader(context)
			hdr, err := tr.Next()

			So(hdr.Name, ShouldEqual, artifactFileName)
			trStr, _ := StreamToString(tr)
			So(trStr, ShouldEqual, "Test artifact file content")

			hdr, err = tr.Next()

			So(hdr.Name, ShouldEqual, "Dockerfile")
			trStr, _ = StreamToString(tr)
			So(trStr, ShouldEqual, testDockerfileContent)
		})
	})
}
