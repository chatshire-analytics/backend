import argparse

from Chatbot import Chatbot


def main(access_token, prompt):
    chatbot = Chatbot(config={
        "access_token": access_token
    })

    # prompt = "I want to you to act like the expert who are good at writing SQL query.\n"
    # prompt += "Given the table below, please write the SQL query that can get the information about Find the most expensive gas fees-cost transaction in last 7 days."
    # prompt += "GIVE A SQL QUERY ONLY **without** any further responses/explanations. If you believe that it is impossible to get the table below, say *NO* without further responses/explanations."
    # prompt += "Given the table below."
    # prompt += "table name: ethereum.core.fact_transactions"
    # prompt += "table schema: BLOCK_HASH text, BLOCK_NUMBER number, BLOCK_TIMESTAMP timestamp_ntz, CUMULATIVE_GAS_USED float, ETH_VALUE float, FROM_ADDRESS text, GAS_LIMIT number, GAS_PRICE float, GAS_USED float"

    response = ""

    for data in chatbot.ask(prompt):
        response = data["message"]

    print(response)

    return response


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("access_token", help="access token for the chatbot")
    parser.add_argument("prompt", help="prompt to pass to the chatbot")
    args = parser.parse_args()

    response = main(args.access_token, args.prompt)
