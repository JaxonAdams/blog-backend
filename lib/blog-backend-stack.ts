import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { Construct } from 'constructs';
import { LambdaFactory, ProjectLambdas } from './lambda/LambdaFactory';
import { APIGatewayFactory } from './apigateway/APIGatewayFactory';
import { S3Factory } from './s3/S3Factory';

export class BlogBackendStack extends cdk.Stack {
    public lambdas: ProjectLambdas;
    public postsBucket: s3.Bucket;

    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // S3 for storing posts
        const s3Factory = new S3Factory(this);
        this.postsBucket = s3Factory.getPostsBucket();

        // Lambdas for API functionality
        const lambdaFactory = new LambdaFactory(this)
        this.lambdas = lambdaFactory.getLambdas();

        // API Gateway for exposing lambdas
        new APIGatewayFactory(this);

        // Misc. tasks now that resources are defined
        s3Factory.grantPermissions();
    }
}
