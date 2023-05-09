resource "null_resource" "function_binary" {
  triggers = {
    always_run = "${timestamp()}"
  }

  provisioner "local-exec" {
    command = "go build -o ${local.binary_path} ${local.src_path}"
    environment = {
      GOOS   = "linux"
      GOARCH = "amd64"
    }
  }
}

data "archive_file" "function_archive" {
  depends_on  = [null_resource.function_binary]
  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

resource "aws_lambda_function" "function" {
  function_name = "bot-smiles"
  description   = "Smiles bot"
  role          = aws_iam_role.lambda.arn
  handler       = local.binary_name
  memory_size   = 128

  filename         = local.archive_path
  source_code_hash = data.archive_file.function_archive.output_base64sha256

  runtime = "go1.x"

  environment {
    variables = {
      PROMOTIONS_URL  = "https://smiles.com.br/promocao"
      USER_AGENT      = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
      EMAIL_TO_NOTIFY = "otaviotech@gmail.com"
      REGEXES         = "^Transferir\\sPontos$"
    }
  }
}

resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.function.function_name}"
  retention_in_days = 7
}
