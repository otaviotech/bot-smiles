resource "aws_cloudwatch_event_rule" "lambda_trigger" {
  name                = "everyday_at_10am_and_16pm"
  description         = "run everyday at 10am and 16pm (Brasilia time)"
  schedule_expression = "cron(0 13,21 * * ? *)"
}

resource "aws_cloudwatch_event_target" "lambda_target" {
  arn  = aws_lambda_function.function.arn
  rule = aws_cloudwatch_event_rule.lambda_trigger.name
}

resource "aws_lambda_permission" "allow_trigger_execution_from_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.lambda_trigger.arn
}
