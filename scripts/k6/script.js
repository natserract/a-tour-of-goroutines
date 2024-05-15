/* eslint-disable no-undef */
import { sleep } from "k6";
import { config } from "./config.js";
import {
  clone,
  generateRandomNumber,
  generateRandomName,
  generateRandomPassword,
} from "./helpers/generator.js";
import { TestProductManagementPost } from "./testCases/productManagement.js";
import exec from "k6/execution";

const stages = [];

if (config.LOAD_TEST) {
  stages.push(
    { target: 50, iterations: 1, duration: "5s" },
    { target: 100, iterations: 1, duration: "10s" },
    { target: 150, iterations: 1, duration: "20s" },
    { target: 200, iterations: 1, duration: "20s" },
    { target: 250, iterations: 1, duration: "20s" },
    { target: 300, iterations: 1, duration: "20s" },
    { target: 600, iterations: 1, duration: "20s" },
    { target: 1200, iterations: 1, duration: "20s" },
  );
} else {
  stages.push({
    target: 1,
    iterations: 1,
  });
}

function determineStage() {
  let elapsedTime = (exec.instance.currentTestRunDuration / 1000).toFixed(0);
  if (elapsedTime < 5) return 1; // First 5 seconds
  if (elapsedTime < 15) return 2; // Next 10 seconds
  if (elapsedTime < 35) return 3; // Next 20 seconds
  if (elapsedTime < 55) return 4; // Next 20 secondsd
  return 5; // Remaining time
}

export const options = {
  stages,
  // summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)'],
};

const positiveCaseConfig = Object.assign(clone(config), {
  POSITIVE_CASE: true,
});

function calculatePercentage(percentage, __VU) {
  return (__VU - 1) % Math.ceil(__VU / Math.round(__VU * percentage)) === 0;
}

const users = [];
function getRandomUser() {
  const i = generateRandomNumber(0, users.length - 1);
  return users[i];
}

export default function () {
  let tags = {};

  if (config.LOAD_TEST) {
    if (determineStage() == 1) {
      let user = {
        name: generateRandomName(),
        password: generateRandomPassword(),
        accessToken:
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUzODg2ODgsInN1YiI6IkFsZmluIiwidWlkIjoiZDNkOWY2NzMtOGRhMy00OWFiLWIyZTktYzVlOTQ2Zjc5ZDQ5In0.jF9wfqyDSzngdMFvB3QzdeZYRtxML_hYPHVpd3yyDQ4",
      };
      users.push(user);
      for (let i = 0; i < 10; i++) {
        TestProductManagementPost(getRandomUser(), positiveCaseConfig, tags);
      }
    } else if (determineStage() == 2) {
      let user = {
        name: generateRandomName(),
        password: generateRandomPassword(),
        accessToken:
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUzODg2ODgsInN1YiI6IkFsZmluIiwidWlkIjoiZDNkOWY2NzMtOGRhMy00OWFiLWIyZTktYzVlOTQ2Zjc5ZDQ5In0.jF9wfqyDSzngdMFvB3QzdeZYRtxML_hYPHVpd3yyDQ4",
      };
      users.push(user);
      for (let i = 0; i < 10; i++) {
        TestProductManagementPost(getRandomUser(), positiveCaseConfig, tags);
      }
    } else if (determineStage() == 3) {
      let user = {
        name: generateRandomName(),
        password: generateRandomPassword(),
        accessToken:
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUzODg2ODgsInN1YiI6IkFsZmluIiwidWlkIjoiZDNkOWY2NzMtOGRhMy00OWFiLWIyZTktYzVlOTQ2Zjc5ZDQ5In0.jF9wfqyDSzngdMFvB3QzdeZYRtxML_hYPHVpd3yyDQ4",
      };
      users.push(user);
      for (let i = 0; i < 10; i++) {
        if (calculatePercentage(0.2, __VU)) {
          TestProductManagementPost(getRandomUser(), positiveCaseConfig, tags);
        }
      }
    }
  } else {
    let user = {
      name: generateRandomName(),
      password: generateRandomPassword(),
      accessToken:
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUzODg2ODgsInN1YiI6IkFsZmluIiwidWlkIjoiZDNkOWY2NzMtOGRhMy00OWFiLWIyZTktYzVlOTQ2Zjc5ZDQ5In0.jF9wfqyDSzngdMFvB3QzdeZYRtxML_hYPHVpd3yyDQ4",
    };
    TestProductManagementPost(user, config, tags);
  }

  sleep(1);
}
