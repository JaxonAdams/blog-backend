import * as cdk from "aws-cdk-lib";
import * as s3 from "aws-cdk-lib/aws-s3";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as authorizers from "aws-cdk-lib/aws-apigatewayv2-authorizers";
import { Construct } from "constructs";
import { LambdaFactory, ProjectLambdas } from "./lambda/LambdaFactory";
import { APIGatewayFactory } from "./apigateway/APIGatewayFactory";
import { S3Factory } from "./s3/S3Factory";
import { DynamoDBFactory } from "./dynamodb/DynamoDBFactory";

export class BlogBackendStack extends cdk.Stack {
  public authorizer: authorizers.HttpLambdaAuthorizer;
  public lambdas: ProjectLambdas;
  public bucket: s3.Bucket;
  public postTable: dynamodb.TableV2;
  public authTable: dynamodb.TableV2;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // S3 for storing posts
    const s3Factory = new S3Factory(this);
    this.bucket = s3Factory.getBucket();

    // DynamoDB for storing post metadata
    const dynamodbFactory = new DynamoDBFactory(this);
    this.postTable = dynamodbFactory.getPostTable();
    this.authTable = dynamodbFactory.getAuthTable();

    // Lambdas for API functionality
    const lambdaFactory = new LambdaFactory(this);
    this.authorizer = lambdaFactory.getAuthorizer();
    this.lambdas = lambdaFactory.getLambdas();

    // API Gateway for exposing lambdas
    new APIGatewayFactory(this);

    this.grantPermissions({
      s3Factory: s3Factory,
      dynamodbFactory: dynamodbFactory,
    });
  }

  private grantPermissions(factories: { [key: string]: any }): void {
    factories["s3Factory"].grantPermissions();
    factories["dynamodbFactory"].grantPermissions();
  }
}
