/* eslint-disable no-loss-of-precision */
import { fail } from "k6";
import {
  MaxInt,
  generateRandomNumber,
  generateRandomName,
  generateTestObjects,
} from "../helpers/generator.js";
import {
  testDeleteAssert,
  testGetAssert,
  testPostJsonAssert,
  testPutJsonAssert,
} from "../helpers/request.js";
import { isUserValid } from "../types/user.js";
import { generateProduct, getRandomCategory } from "../types/product.js";
import {
  isEqualWith,
  isExists,
  isOrdered,
  isTotalDataInRange,
  isValidDate,
  isValidUrl,
} from "../helpers/assertion.js";

const productNegativePayload = (positivePayload) =>
  generateTestObjects(
    {
      name: { notNull: true, type: "string", minLength: 1, maxLength: 30 },
      sku: { notNull: true, type: "string", minLength: 1, maxLength: 30 },
      category: {
        notNull: true,
        type: "string",
        enum: ["Clothing", "Accessories", "Footwear", "Beverages"],
      },
      imageUrl: { notNull: true, type: "string", isUrl: true },
      notes: { notNull: true, type: "string", minLength: 1, maxLength: 200 },
      price: { notNull: true, type: "number", min: 1 },
      stock: { notNull: true, type: "number", min: -1, max: 100000 },
      location: { notNull: true, type: "string", minLength: 1, maxLength: 200 },
      isAvailable: { notNull: true, type: "boolean" },
    },
    positivePayload,
  );

/**
 *
 * @param {import("../config.js").Config} config
 * @param {Object} tags
 * @param {import("../types/user.js").User} user
 * @returns {import("../types/product.js").Product}
 */
export function TestProductManagementPost(user, config, tags) {
  const currentRoute = `${config.BASE_URL}/v1/product`;
  const currentFeature = "post product";

  if (!isUserValid(user)) {
    fail(`${currentFeature} Invalid user object`);
  }

  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };

  /** @type {import("../types/product.js").Product} */
  const productPositivePayload = generateProduct();

  /** @type {import("../helpers/request.js").RequestAssertResponse} */
  let res;

  if (!config.POSITIVE_CASE) {
    productNegativePayload(productPositivePayload).forEach((payload) => {
      // console.log("INVALID payload", payload);

      testPostJsonAssert(
        currentFeature,
        "invalid payload",
        currentRoute,
        payload,
        headers,
        {
          ["should return 400"]: (res) => res.status === 400,
        },
        config,
        tags,
      );
    });
  }

  res = testPostJsonAssert(
    currentFeature,
    "add product with correct payload",
    currentRoute,
    productPositivePayload,
    headers,
    {
      ["should return 201"]: (res) => res.status === 201,
      ["should have id"]: (res) => isExists(res, "data.id"),
      ["should have createdAt and format should be date"]: (res) =>
        isEqualWith(res, "data.createdAt", (v) =>
          v.every((a) => isValidDate(a)),
        ),
    },
    config,
    tags,
  );

  if (res.isSuccess) {
    return Object.assign(productPositivePayload, {
      id: res.res.json().data.id,
      createdAt: res.res.json().data.createdAt,
    });
  }
  return null;
}
