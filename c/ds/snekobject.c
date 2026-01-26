#include "snekobject.h"
#include "assert.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void refcount_free(snek_object_t *obj) {
  switch (obj->kind) {
  case INTEGER:
  case FLOAT:
    break;
  case STRING:
    free(obj->data.v_string);
    break;
  case VECTOR3: {
    refcount_dec(obj->data.v_vector3.x);
    refcount_dec(obj->data.v_vector3.y);
    refcount_dec(obj->data.v_vector3.z);
    break;
  }
  case ARRAY: {
    for (int i = 0; i < obj->data.v_array.size; i++) {
      refcount_dec(obj->data.v_array.elements[i]);
    }
    free(obj->data.v_array.elements);
    break;
  }
  default:
    assert(false);
  }

  free(obj);
}

void refcount_inc(snek_object_t *obj) {
  if (obj == NULL) {
    return;
  }
  obj->refcount++;
}

void refcount_dec(snek_object_t *obj) {
  if (obj == NULL) {
    return;
  }
  obj->refcount--;
  if (obj->refcount == 0) {
    refcount_free(obj);
  }
}

snek_object_t *_new_snek_object() {
  snek_object_t *obj = calloc(1, sizeof(snek_object_t));
  if (obj == NULL) {
    return NULL;
  }
  obj->refcount = 1;
  return obj;
}

bool snek_array_set(snek_object_t *snek_obj, size_t index,
                    snek_object_t *value) {
  if (snek_obj == NULL || value == NULL || snek_obj->kind != ARRAY ||
      snek_obj->data.v_array.size <= index) {
    return false;
  }

  // Decrement old element's refcount if replacing
  if (snek_obj->data.v_array.elements[index] != NULL) {
    refcount_dec(snek_obj->data.v_array.elements[index]);
  }
  // Store new element and increment its refcount
  snek_obj->data.v_array.elements[index] = value;
  refcount_inc(value);
  return true;
}

snek_object_t *snek_array_get(snek_object_t *snek_obj, size_t index) {
  if (snek_obj == NULL || snek_obj->kind != ARRAY ||
      snek_obj->data.v_array.size <= index) {
    return NULL;
  }
  return snek_obj->data.v_array.elements[index];
}

int snek_length(snek_object_t *obj) {
  if (obj == NULL) {
    return -1;
  }

  switch (obj->kind) {
  case INTEGER:
    return 1;
  case FLOAT:
    return 1;
  case STRING:
    return strlen(obj->data.v_string);
  case VECTOR3:
    return 3;
  case ARRAY:
    return obj->data.v_array.size;
  default:
    return -1;
  }
}

snek_object_t *snek_add(snek_object_t *a, snek_object_t *b) {
  if (a == NULL || b == NULL) {
    return NULL;
  }

  switch (a->kind) {
  case INTEGER:
    if (b->kind == a->kind) {
      return new_snek_integer(a->data.v_int + b->data.v_int);
    } else if (b->kind == FLOAT) {
      return new_snek_float(a->data.v_int + b->data.v_float);
    } else {
      return NULL;
    }
  case FLOAT:
    if (b->kind == a->kind) {
      return new_snek_float(a->data.v_float + b->data.v_float);
    } else if (b->kind == INTEGER) {
      return new_snek_integer(a->data.v_float + b->data.v_int);
    } else {
      return NULL;
    }
  case STRING:
    if (b->kind != a->kind) {
      return NULL;
    }
    int s_a = strlen(a->data.v_string);
    int s_b = strlen(b->data.v_string);
    int total_len = s_a + s_b;
    char *temp_str = calloc(total_len + 1, sizeof(char));
    strcat(temp_str, a->data.v_string);
    strcat(temp_str, b->data.v_string);
    snek_object_t *obj = new_snek_string(temp_str);
    free(temp_str);
    return obj;
  case VECTOR3:
    if (b->kind != a->kind) {
      return NULL;
    }
    snek_object_t *x = snek_add(a->data.v_vector3.x, b->data.v_vector3.x);
    snek_object_t *y = snek_add(a->data.v_vector3.y, b->data.v_vector3.y);
    snek_object_t *z = snek_add(a->data.v_vector3.z, b->data.v_vector3.z);
    return new_snek_vector3(x, y, z);
  case ARRAY:
    if (b->kind != a->kind) {
      return NULL;
    }
    int a_len = snek_length(a);
    int b_len = snek_length(b);
    snek_object_t *ab_arr = new_snek_array(a_len + b_len);
    // copy a
    for (int i = 0; i < a_len; i++) {
      snek_object_t *a_elem = snek_array_get(a, i);
      snek_array_set(ab_arr, i, a_elem);
    }
    // copy b after a i.e a_len + i  = offset for b's set
    for (int i = 0; i < b_len; i++) {
      snek_object_t *b_elem = snek_array_get(b, i);
      snek_array_set(ab_arr, i + a_len, b_elem); // NOTE:offset
    }
    return ab_arr;
  default:
    return NULL;
  }
}

snek_object_t *new_snek_array(size_t size) {
  snek_object_t *obj = malloc(sizeof(snek_object_t));
  if (obj == NULL) {
    return NULL;
  }
  obj->refcount = 1;

  // init the pointers to the array
  snek_object_t **array = calloc(size, sizeof(snek_object_t *));
  obj->kind = ARRAY;
  snek_array_t array_data = {
      .size = size,
      .elements = array,
  };
  obj->data.v_array = array_data;
  return obj;
}

snek_object_t *new_snek_vector3(snek_object_t *x, snek_object_t *y,
                                snek_object_t *z) {
  if (x == NULL || y == NULL || z == NULL) {
    return NULL;
  }

  // NOTE: we use size of for types, else use given size/length
  snek_object_t *obj = malloc(sizeof(snek_object_t));
  if (obj == NULL) {
    return NULL;
  }

  obj->kind = VECTOR3;
  snek_vector_t vector3 = {.x = x, .y = y, .z = z};
  obj->data.v_vector3 = vector3;
  obj->refcount = 1;
  refcount_inc(x);
  refcount_inc(y);
  refcount_inc(z);
  return obj;
}

snek_object_t *new_snek_integer(int value) {
  snek_object_t *obj = malloc(sizeof(snek_object_t));
  if (obj == NULL) {
    return NULL;
  }
  obj->refcount = 1;

  obj->kind = INTEGER;
  obj->data.v_int = value;
  return obj;
}

snek_object_t *new_snek_float(float value) {
  snek_object_t *obj = malloc(sizeof(snek_object_t));
  if (obj == NULL) {
    return NULL;
  }
  obj->refcount = 1;

  obj->kind = FLOAT;
  obj->data.v_float = value;
  return obj;
}

snek_object_t *new_snek_string(char *value) {
  snek_object_t *obj = malloc(sizeof(snek_object_t));
  if (obj == NULL) {
    return NULL;
  }
  obj->refcount = 1;

  int len = strlen(value);
  // NOTE: we have a size already and this isn't a type
  char *dst = malloc(len + 1);
  if (dst == NULL) {
    free(obj);
    return NULL;
  }

  strcpy(dst, value);

  obj->kind = STRING;
  obj->data.v_string = dst;
  return obj;
}
