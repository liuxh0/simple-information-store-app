const aws = require('aws-sdk');

/**
 *
 * Event doc: https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html#api-gateway-simple-proxy-for-lambda-input-format
 * @param {Object} event - API Gateway Lambda Proxy Input Format
 *
 * Context doc: https://docs.aws.amazon.com/lambda/latest/dg/nodejs-prog-model-context.html
 * @param {Object} context
 *
 * Return doc: https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
 * @returns {Object} object - API Gateway Lambda Proxy Output Format
 *
 */
exports.lambdaHandler = async (event, context) => {
  const id = event.pathParameters.id;

  const ddb = new aws.DynamoDB(getDynamoDbOption());
  const item = await ddb.getItem({
    TableName: getValuesTableName(),
    Key:{
      'Id': { S: id }
    }
  }).promise();

  if (item.Item == undefined) {
    return {
      'statusCode': 404
    };
  }

  const value = item.Item['Value'].S;
  return {
    'statusCode': 200,
    'body': value
  };
};

function getDynamoDbOption() {
  if (process.env.AWS_SAM_LOCAL) {
    return {
      endpoint: 'http://docker.for.mac.localhost:8000'
    };
  } else {
    return undefined;
  }
}

function getValuesTableName() {
  if (process.env.AWS_SAM_LOCAL) {
    return 'Values';
  } else {
    return process.env.VALUES_TABLE_NAME;
  }
}
