import * as cdk from "aws-cdk-lib";
import * as s3 from "aws-cdk-lib/aws-s3";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";
import { LambdaFactory, ProjectLambdas } from "./lambda/LambdaFactory";
import { APIGatewayFactory } from "./apigateway/APIGatewayFactory";
import { S3Factory } from "./s3/S3Factory";
import { DynamoDBFactory } from "./apigateway/dynamodb/DynamoDBFactory";

export class BlogBackendStack extends cdk.Stack {
  public lambdas: ProjectLambdas;
  public bucket: s3.Bucket;
  public table: dynamodb.TableV2;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // S3 for storing posts
    const s3Factory = new S3Factory(this);
    this.bucket = s3Factory.getBucket();

    // DynamoDB for storing post metadata
    const dynamodbFactory = new DynamoDBFactory(this);
    this.table = dynamodbFactory.getTable();

    // Lambdas for API functionality
    const lambdaFactory = new LambdaFactory(this);
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
