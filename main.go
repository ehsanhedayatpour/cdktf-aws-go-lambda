package main

import (
	"fmt"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"

	"cdk.tf/go/stack/generated/hashicorp/aws"
	data "cdk.tf/go/stack/generated/hashicorp/aws/datasources"
	"cdk.tf/go/stack/generated/hashicorp/aws/iam"
	lambda "cdk.tf/go/stack/generated/hashicorp/aws/lambdafunction"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	// AWS region
	region := cdktf.NewTerraformVariable(stack, jsii.String("AWS_REGION"), &cdktf.TerraformVariableConfig{
		Default:     "us-east-1",
		Description: jsii.String("Choose which AWS region to use"),
		Nullable:    jsii.Bool(false),
		Sensitive:   jsii.Bool(false),
		Type:        jsii.String("string"),
	})
	aws.NewAwsProvider(stack, jsii.String("aws"), &aws.AwsProviderConfig{
		Region: region.StringValue(),
	})

	// IAM Role for Lambda function
	assumeRolePolicyJson := `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"Service": "lambda.amazonaws.com"},
			"Action": "sts:AssumeRole"
		}]
	}`
	iamRole := iam.NewIamRole(stack, jsii.String("cdktf-go-example-function-role"), &iam.IamRoleConfig{
		Name:              jsii.String("cdktf-example-function-role"),
		AssumeRolePolicy:  jsii.String(assumeRolePolicyJson),
		ManagedPolicyArns: jsii.Strings("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"),
	})

	// Create Lambda function
	self := data.NewDataAwsCallerIdentity(stack, jsii.String("caller-identity"), &data.DataAwsCallerIdentityConfig{})
	accountId := self.AccountId()
	imageUri := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/go-hello-world:latest", *accountId, *region.StringValue())
	arch := cdktf.NewTerraformVariable(stack, jsii.String("CONTAINER_IMG_ARCH"), &cdktf.TerraformVariableConfig{
		Default:     "x86_64",
		Description: jsii.String("The valid value is 'x86_64' or 'arm64'."),
		Nullable:    jsii.Bool(false),
		Sensitive:   jsii.Bool(false),
		Type:        jsii.String("string"),
	})
	lambda.NewLambdaFunction(stack, jsii.String("cdktf-go-example-function"), &lambda.LambdaFunctionConfig{
		FunctionName:  jsii.String("cdktf-example-function"),
		Role:          jsii.String(*iamRole.Arn()),
		Architectures: jsii.Strings(*arch.StringValue()),
		PackageType:   jsii.String("Image"),
		ImageUri:      jsii.String(imageUri),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)

	NewMyStack(app, "cdktf-go-example")
	app.Synth()
}
