import * as cdk from 'aws-cdk-lib';
import * as aws_apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import { ProjectLambdas } from '../lambda/LambdaFactory';

export class APIGatewayFactory {
    private stack: cdk.Stack;
    private gateway: aws_apigatewayv2.HttpApi;

    constructor(stack: cdk.Stack, lambdas: ProjectLambdas) {
        this.stack = stack;
        this.gateway = this.makeHttpApi();

        this.loadRoutes(lambdas);
        this.makeCfnOutputs();
    }

    private makeHttpApi(): aws_apigatewayv2.HttpApi {
        return new aws_apigatewayv2.HttpApi(this.stack, 'HttpApi', {
            apiName: this.stack.stackName,
        });
    }

    private makeCfnOutputs(): void {
        new cdk.CfnOutput(this.stack, 'HttpApiUrlReference', {
            exportName: `${this.stack.stackName}-HttpApiUrl`,
            value: this.gateway.apiEndpoint,
            description: 'HTTP API URL',
        });
    }

    private loadRoutes(lambdas: ProjectLambdas): void {
        const {
            createPostLambda,
        } = lambdas;

        this.gateway.addRoutes({
            path: '/api/v1/posts',
            methods: [aws_apigatewayv2.HttpMethod.POST],
            integration: new cdk.aws_apigatewayv2_integrations.HttpLambdaIntegration(
                'CreatePostIntegration',
                createPostLambda,
            )
        });
    }
}
