package OctaForce

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"io/ioutil"
	"log"
	"strings"
)

func compileShader(path string, shaderType uint32) (uint32, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	shader := gl.CreateShader(shaderType)

	source, free := gl.Strs(string(content))
	gl.ShaderSource(shader, 1, source, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
