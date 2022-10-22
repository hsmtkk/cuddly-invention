import { Construct } from "constructs";
import { App, TerraformStack, CloudBackend, NamedCloudWorkspace, TerraformAsset, AssetType } from "cdktf";
import * as google from '@cdktf/provider-google';
import * as path from 'path';
import { DataGoogleIamPolicy } from "@cdktf/provider-google";

const default_location = 'asia-northeast1';
const project_id = 'cuddly-invention';

class MyStack extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);

    new google.GoogleProvider(this, 'Google', {
      project: project_id,
    });

    const asset_bucket = new google.StorageBucket(this, 'asset_bucket', {
      location: default_location,
      name: `asset_bucket_${project_id}`,      
    });

    // To allow unauthenticated access, I have to allow unauthenticated access at the underlying Cloud Run service.
    const noauth_policy = new DataGoogleIamPolicy(this, 'noauth_policy', {
      binding: [{
        role: 'roles/run.invoker',
        members: ['allUsers'],
      }],
    });

    const cookietest_zip = new TerraformAsset(this, 'cookietest_zip', {
      path: path.resolve('cookietest'),
      type: AssetType.ARCHIVE,
    });

    const cookietest_object = new google.StorageBucketObject(this, 'cookietest_object', {
      bucket: asset_bucket.name,
      name: cookietest_zip.assetHash,
      source: cookietest_zip.path,
    });

    const cookietest_func = new google.Cloudfunctions2Function(this, 'cookietest', {
      buildConfig: {
        entryPoint: 'cookietest',
        runtime: 'go119',
        source: {
          storageSource: {
            bucket: asset_bucket.name,
            object: cookietest_object.name,
          },
        }
      },
      location: default_location,
      name: 'cookietest',
    });

    new google.CloudRunServiceIamPolicy(this, 'cookietest_noauth', {
      location: default_location,
      policyData: noauth_policy.policyData,
      service: cookietest_func.name,
    });
  }
}

const app = new App();
const stack = new MyStack(app, "cuddly-invention");
new CloudBackend(stack, {
  hostname: "app.terraform.io",
  organization: "hsmtkkdefault",
  workspaces: new NamedCloudWorkspace("cuddly-invention")
});
app.synth();
