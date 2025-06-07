#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { BlogBackendStack } from "../lib/blog-backend-stack";

const app = new cdk.App();
new BlogBackendStack(app, "BlogBackendStack", {});

