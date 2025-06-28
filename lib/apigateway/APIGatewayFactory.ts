import * as cdk from "aws-cdk-lib";
import * as aws_apigatewayv2 from "aws-cdk-lib/aws-apigatewayv2";
import { BlogBackendStack } from "../blog-backend-stack";

export class APIGatewayFactory {
  private stack: BlogBackendStack;
  private gateway: aws_apigatewayv2.HttpApi;

  constructor(stack: BlogBackendStack) {
    this.stack = stack;
    this.gateway = this.makeHttpApi();

    this.loadRoutes();
    this.makeCfnOutputs();
  }

  private makeHttpApi(): aws_apigatewayv2.HttpApi {
    return new aws_apigatewayv2.HttpApi(this.stack, "HttpApi", {
      apiName: this.stack.stackName,
      corsPreflight: {
        allowHeaders: ["Content-Type", "Authorization"],
        allowMethods: [
          aws_apigatewayv2.CorsHttpMethod.GET,
          aws_apigatewayv2.CorsHttpMethod.POST,
          aws_apigatewayv2.CorsHttpMethod.PATCH,
          aws_apigatewayv2.CorsHttpMethod.DELETE,
          aws_apigatewayv2.CorsHttpMethod.OPTIONS,
        ],
        allowOrigins: ["*"], // TODO: restrict me
        maxAge: cdk.Duration.days(10),
      },
    });
  }

  private makeCfnOutputs(): void {
    new cdk.CfnOutput(this.stack, "HttpApiUrlReference", {
      exportName: `${this.stack.stackName}-HttpApiUrl`,
      value: this.gateway.apiEndpoint,
      description: "HTTP API URL",
    });
  }

  private loadRoutes(): void {
    const authorizer = this.stack.authorizer;
    const {
      createPostLambda,
      updatePostLambda,
      getPostByIdLambda,
      getAllPostsLambda,
      deletePostLambda,
      loginAdminLambda,
    } = this.stack.lambdas;

    this.gateway.addRoutes({
      path: "/api/v1/posts",
      methods: [aws_apigatewayv2.HttpMethod.POST],
      integration: new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
        "CreatePostIntegration",
        createPostLambda,
      ),
      authorizer,
    });

    this.gateway.addRoutes({
      path: "/api/v1/posts",
      methods: [aws_apigatewayv2.HttpMethod.GET],
      integration: new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
        "GetAllPostsIntegration",
        getAllPostsLambda,
      ),
    });

    this.gateway.addRoutes({
      path: "/api/v1/posts/{post_id}",
      methods: [aws_apigatewayv2.HttpMethod.PATCH],
      integration: new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
        "UpdatePostIntegration",
        updatePostLambda,
      ),
      authorizer,
    });

    this.gateway.addRoutes({
      path: "/api/v1/posts/{post_id}",
      methods: [aws_apigatewayv2.HttpMethod.GET],
      integration: new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
        "GetPostByIdIntegration",
        getPostByIdLambda,
      ),
    });

    this.gateway.addRoutes({
      path: "/api/v1/posts/{post_id}",
      methods: [aws_apigatewayv2.HttpMethod.DELETE],
      integration: new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
        "DeletePostIntegration",
        deletePostLambda,
      ),
      authorizer,
    });

    this.gateway.addRoutes({
      path: "/api/v1/auth/login/admin",
      methods: [aws_apigatewayv2.HttpMethod.POST],
      integration: new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
        "LoginAdminIntegration",
        loginAdminLambda,
      ),
    });
  }
}
