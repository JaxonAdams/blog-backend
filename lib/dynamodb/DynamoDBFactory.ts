import * as cdk from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { BlogBackendStack } from "../blog-backend-stack";

export class DynamoDBFactory {
  private stack: BlogBackendStack;
  private postTable: dynamodb.TableV2;
  private authTable: dynamodb.TableV2;

  constructor(stack: BlogBackendStack) {
    this.stack = stack;
    this.postTable = this.makePostTable();
    this.authTable = this.makeAuthTable();

    this.makeCfnOutputs();
  }

  private makePostTable(): dynamodb.TableV2 {
    return new dynamodb.TableV2(this.stack, "PostMetadataTable", {
      tableName: `${this.stack.stackName}-PostMetadataTable`,
      partitionKey: { name: "id", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "createdAt", type: dynamodb.AttributeType.NUMBER },
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });
  }

  private makeAuthTable(): dynamodb.TableV2 {
    return new dynamodb.TableV2(this.stack, "AuthTable", {
      tableName: `${this.stack.stackName}-AuthTable`,
      partitionKey: { name: "username", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "modifiedAt", type: dynamodb.AttributeType.NUMBER },
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });
  }

  private makeCfnOutputs(): void {
    new cdk.CfnOutput(this.stack, "PostMetadataTableNameReference", {
      exportName: `${this.stack.stackName}-PostMetadataTableName`,
      value: this.postTable.tableName,
      description: "DynamoDB Table Name",
    });
    new cdk.CfnOutput(this.stack, "PostMetadataTableARNReference", {
      exportName: `${this.stack.stackName}-PostMetadataTableARN`,
      value: this.postTable.tableArn,
      description: "DynamoDB Table ARN",
    });
  }

  public grantPermissions(): void {
    const {
      createPostLambda,
      updatePostLambda,
      getPostByIdLambda,
      getAllPostsLambda,
      deletePostLambda,
      loginAdminLambda,
    } = this.stack.lambdas;

    this.postTable.grantWriteData(createPostLambda);

    this.postTable.grantReadData(getPostByIdLambda);
    this.postTable.grantReadData(getAllPostsLambda);

    this.postTable.grantReadWriteData(updatePostLambda);
    this.postTable.grantReadWriteData(deletePostLambda);

    this.authTable.grantReadData(loginAdminLambda);
  }

  public getPostTable(): dynamodb.TableV2 {
    return this.postTable;
  }

  public getAuthTable(): dynamodb.TableV2 {
    return this.authTable;
  }
}
