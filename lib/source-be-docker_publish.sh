function docker_publish_jenkins_context {
    return
}

function docker_publish_check_and_set {
    return
}

function docker_publish_set_path {
    return
}

function unset_docker_publish {
    return
}

# Used to create wrapper files.
# List of files are defined at bottom in beWrappers
function be_create_wrapper_docker_publish {
    case $1 in
        "publish_alltags.sh") 
            cat $BASE_DIR/modules/docker_publish/bin/$1 >> $2
            ;;
    esac
}

function be_docker_publish_mount_setup {
    return
}

function be_do_docker_publish_docker_run {
    return
}

function be_create_docker_publish_docker_build {
    return
}

function docker_publish_create_build_env {
    return
}

beWrappers["docker_publish"]="publish_alltags.sh"