#version 330
uniform mat4 projection;
uniform mat4 camera;
in vec3 vertexPosition;

in vec3 instanceColor;
in vec4 transformX;
in vec4 transformY;
in vec4 transformZ;
in vec4 transformS;

out vec3 color;

void main() {
    color = instanceColor;
    gl_Position =
        projection *
        camera *
        mat4(transformX, transformY, transformZ, transformS) *
        vec4(vertexPosition, 1);
}
