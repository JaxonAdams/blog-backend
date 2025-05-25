import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { BlogBackendStack } from '../blog-backend-stack';

export class S3Factory {
    private stack: BlogBackendStack;
    private postsBucket: s3.Bucket;

    constructor(stack: BlogBackendStack) {
        this.stack = stack;
        this.postsBucket = this.makePostsBucket();

        this.makeCfnOutputs();
    }

    private makePostsBucket(): s3.Bucket {
        return new s3.Bucket(this.stack, 'Posts', {
            bucketName: `${this.stack.stackName}-Posts`.toLowerCase(),
            versioned: true,
        });
    }

    private makeCfnOutputs(): void {
        new cdk.CfnOutput(this.stack, 'PostsBucketNameReference', {
            exportName: `${this.stack.stackName}-PostsBucketName`,
            value: this.postsBucket.bucketName,
            description: 'S3 Bucket Name',
        });
    }

    public grantPermissions(): void {
        const {
            createPostLambda,
        } = this.stack.lambdas;

        this.postsBucket.grantWrite(createPostLambda);
    }

    public getPostsBucket(): s3.Bucket {
        return this.postsBucket;
    }
}
