import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as authorizers from "aws-cdk-lib/aws-apigatewayv2-authorizers";
import { BlogBackendStack } from "../blog-backend-stack";

export class LambdaFactory {
  private stack: BlogBackendStack;
  private authorizer: authorizers.HttpLambdaAuthorizer;
  private lambdas: ProjectLambdas;

  constructor(stack: BlogBackendStack) {
    this.stack = stack;
    this.authorizer = this.makeAuthorizerLambda();

    this.lambdas = {
      createPostLambda: this.makeCreatePostLambda(),
      updatePostLambda: this.makeUpdatePostLambda(),
      getPostByIdLambda: this.makeGetPostByIdLambda(),
      getAllPostsLambda: this.makeGetAllPostsLambda(),
      deletePostLambda: this.makeDeletePostLambda(),
      loginAdminLambda: this.makeLoginAdminLambda(),
    };
  }

  private makeAuthorizerLambda(): authorizers.HttpLambdaAuthorizer {
    const lambdaFn = new lambda.Function(this.stack, "AuthorizerFunction", {
      functionName: `${this.stack.stackName}-Authorizer`,
      runtime: lambda.Runtime.PROVIDED_AL2023,
      timeout: cdk.Duration.seconds(30),
      handler: "main",
      code: lambda.Code.fromAsset("src/api/auth/authorizer/build"),
      environment: {
        JWT_SECRET: process.env.JWT_SECRET || "",
      },
    });

    return new authorizers.HttpLambdaAuthorizer(
      "BlogLambdaAuthorizer",
      lambdaFn,
      {
        responseTypes: [authorizers.HttpLambdaResponseType.SIMPLE],
      },
    );
  }

  private makeCreatePostLambda(): lambda.Function {
    return new lambda.Function(this.stack, "CreatePost", {
      functionName: `${this.stack.stackName}-CreatePost`,
      runtime: lambda.Runtime.PROVIDED_AL2023,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset("src/api/post/create/build"),
      handler: "bootstrap",
      environment: {
        S3_BUCKET_NAME: this.stack.bucket.bucketName,
        POST_METADATA_TABLE_NAME: this.stack.postTable.tableName,
      },
    });
  }

  private makeUpdatePostLambda(): lambda.Function {
    return new lambda.Function(this.stack, "UpdatePost", {
      functionName: `${this.stack.stackName}-UpdatePost`,
      runtime: lambda.Runtime.PROVIDED_AL2023,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset("src/api/post/update/build"),
      handler: "bootstrap",
      environment: {
        S3_BUCKET_NAME: this.stack.bucket.bucketName,
        POST_METADATA_TABLE_NAME: this.stack.postTable.tableName,
      },
    });
  }

  private makeGetPostByIdLambda(): lambda.Function {
    return new lambda.Function(this.stack, "GetPostById", {
      functionName: `${this.stack.stackName}-GetPostById`,
      runtime: lambda.Runtime.PROVIDED_AL2023,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset("src/api/post/getbyid/build"),
      handler: "bootstrap",
      environment: {
        S3_BUCKET_NAME: this.stack.bucket.bucketName,
        S3_URL_EXPIRY_SECONDS: "3600",
        POST_METADATA_TABLE_NAME: this.stack.postTable.tableName,
      },
    });
  }

  private makeGetAllPostsLambda(): lambda.Function {
    return new lambda.Function(this.stack, "GetAllPosts", {
      functionName: `${this.stack.stackName}-GetAllPosts`,
      runtime: lambda.Runtime.PROVIDED_AL2023,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset("src/api/post/getall/build"),
      handler: "bootstrap",
      environment: {
        POST_METADATA_TABLE_NAME: this.stack.postTable.tableName,
        DEFAULT_PAGE_SIZE: "20",
      },
    });
  }

  private makeDeletePostLambda(): lambda.Function {
    return new lambda.Function(this.stack, "DeletePost", {
      functionName: `${this.stack.stackName}-DeletePost`,
      runtime: lambda.Runtime.PROVIDED_AL2023,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset("src/api/post/delete/build"),
      handler: "bootstrap",
      environment: {
        POST_METADATA_TABLE_NAME: this.stack.postTable.tableName,
      },
    });
  }

  private makeLoginAdminLambda(): lambda.Function {
    return new lambda.Function(this.stack, "LoginAdmin", {
      functionName: `${this.stack.stackName}-LoginAdmin`,
      runtime: lambda.Runtime.PROVIDED_AL2023,
      timeout: cdk.Duration.seconds(30),
      code: lambda.Code.fromAsset("src/api/auth/login/admin/build"),
      handler: "bootstrap",
      environment: {
        AUTH_TABLE_NAME: this.stack.authTable.tableName,
        JWT_SECRET: process.env.JWT_SECRET || "",
      },
    });
  }

  public getLambdas(): ProjectLambdas {
    return this.lambdas;
  }

  public getAuthorizer(): authorizers.HttpLambdaAuthorizer {
    return this.authorizer;
  }
}

export type ProjectLambdas = {
  [key: string]: lambda.Function;
};
