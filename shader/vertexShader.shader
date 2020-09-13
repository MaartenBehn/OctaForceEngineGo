#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 transform;
in vec3 vert;

void main() {
    gl_Position = projection * camera * transform * vec4(vert, 1);
}
