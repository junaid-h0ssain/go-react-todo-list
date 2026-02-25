import { mutation, query } from "./_generated/server";
import { v } from "convex/values";

export const list = query({
  args: {},
  handler: async (ctx) => {
    return await ctx.db.query("todos").collect();
  },
});

export const add = mutation({
  args: {
    body: v.string(),
  },
  handler: async (ctx, args) => {
    const id = await ctx.db.insert("todos", {
      body: args.body,
      completed: false,
    });
    return await ctx.db.get(id);
  },
});

export const setCompleted = mutation({
  args: {
    id: v.id("todos"),
  },
  handler: async (ctx, args) => {
    const todo = await ctx.db.get(args.id);
    if (!todo) {
      throw new Error("Todo not found");
    }
    await ctx.db.patch(args.id, { completed: true });
    return await ctx.db.get(args.id);
  },
});

export const remove = mutation({
  args: {
    id: v.id("todos"),
  },
  handler: async (ctx, args) => {
    const todo = await ctx.db.get(args.id);
    if (!todo) {
      throw new Error("Todo not found");
    }
    await ctx.db.delete(args.id);
  },
});
