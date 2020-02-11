resource "rabbitmq_user" "test" {
  name     = "test"
  password = "foobar"
  tags     = ["management"]
}

resource "rabbitmq_permissions" "test" {
  user  = rabbitmq_user.test.name
  vhost = "/"

  permissions {
    write     = ""
    read      = ""
    configure = ""
  }
}

