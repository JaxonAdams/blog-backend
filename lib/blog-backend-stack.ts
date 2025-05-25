import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import { LambdaFactory } from './lambda/LambdaFactory';
import { APIGatewayFactory } from './apigateway/APIGatewayFactory';

export class BlogBackendStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const lambdaFactory = new LambdaFactory(this)
        const lambdas = lambdaFactory.getLambdas();

        new APIGatewayFactory(this, lambdas);

    }
}
