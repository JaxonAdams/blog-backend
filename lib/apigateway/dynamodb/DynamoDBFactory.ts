import * as cdk from 'aws-cdk-lib';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import { BlogBackendStack } from '../../blog-backend-stack';

export class DynamoDBFactory {
    private stack: BlogBackendStack;
    private table: dynamodb.TableV2;

    constructor(stack: BlogBackendStack) {
        this.stack = stack;
        this.table = this.makeTable();

        this.makeCfnOutputs();
    }

    private makeTable(): dynamodb.TableV2 {
        return new dynamodb.TableV2(this.stack, 'PostMetadataTable', {
            tableName: `${this.stack.stackName}-PostMetadataTable`,
            partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
            sortKey: { name: 'createdAt', type: dynamodb.AttributeType.NUMBER },
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        });
    }

    private makeCfnOutputs(): void {
        new cdk.CfnOutput(this.stack, 'PostMetadataTableNameReference', {
            exportName: `${this.stack.stackName}-PostMetadataTableName`,
            value: this.table.tableName,
            description: 'DynamoDB Table Name',
        })
        new cdk.CfnOutput(this.stack, 'PostMetadataTableARNReference', {
            exportName: `${this.stack.stackName}-PostMetadataTableARN`,
            value: this.table.tableArn,
            description: 'DynamoDB Table ARN',
        })
    }

    public grantPermissions(): void {
        const {
            createPostLambda,
            getPostByIdLambda,
        } = this.stack.lambdas;

        this.table.grantWriteData(createPostLambda);
        this.table.grantReadData(getPostByIdLambda);
    }

    public getTable(): dynamodb.TableV2 {
        return this.table;
    }

}
