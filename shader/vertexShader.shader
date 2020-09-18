#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 transform;
in vec3 vertexPosition;
in vec3 vertexColor;
out vec3 color;

void main() {
    color = vertexColor;
    gl_Position =  projection * camera * transform * vec4(vertexPosition, 1);
}
