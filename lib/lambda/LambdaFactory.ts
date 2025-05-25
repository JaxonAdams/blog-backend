import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import { BlogBackendStack } from '../blog-backend-stack';

export class LambdaFactory {
    private stack: BlogBackendStack;
    private lambdas: ProjectLambdas;

    constructor(stack: BlogBackendStack) {
        this.stack = stack;

        this.lambdas = {
            createPostLambda: this.makeCreatePostLambda(),
        };
    }

    private makeCreatePostLambda(): lambda.Function {
        return new lambda.Function(this.stack, 'CreatePost', {
            functionName: `${this.stack.stackName}-CreatePost`,
            runtime: lambda.Runtime.PROVIDED_AL2023,
            timeout: cdk.Duration.seconds(30),
            code: lambda.Code.fromAsset('src/api/post/create/build'),
            handler: 'bootstrap',
            environment: {
                'S3_BUCKET_NAME': this.stack.bucket.bucketName,
            },
        });
    }

    public getLambdas(): ProjectLambdas {
        return this.lambdas;
    }
}

export type ProjectLambdas = {
    [key: string]: lambda.Function
};
