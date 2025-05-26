import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { BlogBackendStack } from '../blog-backend-stack';

export class S3Factory {
    private stack: BlogBackendStack;
    private bucket: s3.Bucket;

    constructor(stack: BlogBackendStack) {
        this.stack = stack;
        this.bucket = this.makeBucket();

        this.makeCfnOutputs();
    }

    private makeBucket(): s3.Bucket {
        return new s3.Bucket(this.stack, 'PostsBucket', {
            bucketName: `${this.stack.stackName}-PostsBucket`.toLowerCase(),
            versioned: true,
        });
    }

    private makeCfnOutputs(): void {
        new cdk.CfnOutput(this.stack, 'PostsBucketNameReference', {
            exportName: `${this.stack.stackName}-PostsBucketName`,
            value: this.bucket.bucketName,
            description: 'S3 Bucket Name',
        });
    }

    public grantPermissions(): void {
        const {
            createPostLambda,
            getPostByIdLambda,
        } = this.stack.lambdas;

        this.bucket.grantWrite(createPostLambda);
        this.bucket.grantRead(getPostByIdLambda);
    }

    public getBucket(): s3.Bucket {
        return this.bucket;
    }
}
