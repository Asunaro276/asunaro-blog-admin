resource "aws_api_gateway_rest_api" "cms" {
  name        = "${var.name}-api"
  description = "API Gateway for Lambda"
}

resource "aws_api_gateway_resource" "cms" {
  rest_api_id = aws_api_gateway_rest_api.cms.id
  parent_id   = aws_api_gateway_rest_api.cms.root_resource_id
  path_part   = "healthcheck"
}

resource "aws_api_gateway_method" "cms" {
  rest_api_id   = aws_api_gateway_rest_api.cms.id
  resource_id   = aws_api_gateway_resource.cms.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "cms" {
  rest_api_id = aws_api_gateway_rest_api.cms.id
  resource_id = aws_api_gateway_resource.cms.id
  http_method = aws_api_gateway_method.cms.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.cms.invoke_arn
}

resource "aws_api_gateway_method_response" "response_200" {
  rest_api_id = aws_api_gateway_rest_api.cms.id
  resource_id = aws_api_gateway_resource.cms.id
  http_method = aws_api_gateway_method.cms.http_method
  status_code = "200"
  response_models = {
    "application/json" = "Empty"
  }
}

resource "aws_lambda_permission" "cms" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cms.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.cms.execution_arn}/*/*"
}

resource "aws_api_gateway_deployment" "cms" {
  rest_api_id = aws_api_gateway_rest_api.cms.id
  depends_on  = [aws_api_gateway_integration.cms]
}

resource "aws_api_gateway_stage" "cms" {
  deployment_id = aws_api_gateway_deployment.cms.id
  rest_api_id   = aws_api_gateway_rest_api.cms.id
  stage_name    = "v1"
}
