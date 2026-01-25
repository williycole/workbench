#include "snekstack.h"
#include <assert.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

void stack_free(stack_t *stack) {
  if (stack == NULL) {
    return;
  }
  free(stack->data);
  free(stack);
}

// stack items go on and off via (last in first out)
void *stack_pop(stack_t *stack) {
  if (stack->count == 0) {
    return NULL;
  }
  stack->count--;
  // NOTE: again, the count is the index of the
  // last item(i.e. where we currently are)
  return stack->data[stack->count];
}

void stack_push(stack_t *stack, void *obj) {
  if (stack->count == stack->capacity) {
    stack->capacity *= 2;

    stack->data = realloc(stack->data, sizeof(void *) * stack->capacity);
    if (stack->data == NULL) {
      stack->capacity /= 2;
      exit(1);
    }
  }

  // NOTE: current count is where we are in the stack currently
  // so it always needs to be the index
  stack->data[stack->count] = obj;
  stack->count++;
}

stack_t *stack_new(size_t capacity) {
  stack_t *stack = malloc(sizeof(stack_t));
  if (stack == NULL) {
    return NULL;
  }

  stack->count = 0;
  stack->capacity = capacity;
  // NOTE: malloc stack, notice realloc is the same pattern but
  //  pointing with the stack thats passed in
  stack->data = malloc(sizeof(void *) * capacity);

  if (stack->data == NULL) {
    free(stack);
    return NULL;
  }

  return stack;
}

void scary_double_push(stack_t *s) {
  stack_push(s, (void *)(int *)1337);
  int *v = malloc(sizeof(int));
  *v = 1024;
  stack_push(s, v);
}

void stack_push_multiple_types(stack_t *s) {
  // create and push float
  float *f = malloc(sizeof(float));
  *f = 3.14;
  stack_push(s, f);

  // create and push char *
  const char *w = "Sneklang is blazingly slow!";
  // malloc len of string plus \0 null term
  char *c = malloc(strlen(w) + 1);
  strcpy(c, w);
  stack_push(s, c);
}
