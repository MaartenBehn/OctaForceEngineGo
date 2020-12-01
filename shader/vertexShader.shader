#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 transform;
uniform vec3 inColor;

in vec3 vertexPosition;

out vec3 outColor;

void main() {
    outColor = inColor;
    gl_Position = projection * camera * transform * vec4(vertexPosition, 1);
}
