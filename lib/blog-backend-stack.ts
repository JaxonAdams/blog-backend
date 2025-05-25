import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as aws_apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import { LambdaFactory } from './lambda/LambdaFactory';

export class BlogBackendStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const lambdaFactory = new LambdaFactory(this)
        const {
            createPostLambda,
        } = lambdaFactory.getLambdas();

        const httpApi = new aws_apigatewayv2.HttpApi(this, 'HttpApi', {
            apiName: this.stackName,
        });

        httpApi.addRoutes({
            path: '/api/v1/posts',
            methods: [aws_apigatewayv2.HttpMethod.POST],
            integration: new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
                'CreatePostIntegration',
                createPostLambda,
            )
        });

        new cdk.CfnOutput(this, 'HttpApiUrlReference', {
            exportName: `${this.stackName}-HttpApiUrl`,
            value: httpApi.apiEndpoint,
            description: 'HTTP API URL',
        });
    }
}
