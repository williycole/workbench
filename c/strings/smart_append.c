#include "smart_append.h"
#include <stdio.h>
#include <string.h>
int smart_append(TextBuffer *dest, const char *src) {
  if (dest == NULL || src == NULL) {
    return 1;
  }

  const size_t BUFFER_CAPACITY = 64;

  size_t src_len = strlen(src);
  // dest_space_left = total_capacity - current_length - 1 (for '\0')
  size_t dest_space_left = BUFFER_CAPACITY - dest->length - 1;

  if (src_len > dest_space_left) {
    printf("partial append\n");

    // not enough room: append what we can and set length
    strncat(dest->buffer, src, dest_space_left);
    dest->length = BUFFER_CAPACITY - 1; // buffer is now full (63 chars)

    return 1;
  } else {
    printf("full append\n");

    // enough room: append entire src and set length
    strcat(dest->buffer, src);
    dest->length = dest->length + src_len;

    return 0;
  }
}

int main(void) {
  // Example 1: enough space (should return 0)
  TextBuffer a;
  strcpy(a.buffer, "Hello");
  a.length = strlen(a.buffer);

  const char *src1 = " World";
  int result1 = smart_append(&a, src1);

  printf("Example 1:\n");
  printf("Src 1: %s\n", src1);
  printf("  result: %d\n", result1);
  printf("  buffer: \"%s\"\n", a.buffer);
  printf("  length: %zu\n\n", a.length);

  // Example 2: not enough space (should return 1, truncates)
  TextBuffer b;
  strcpy(b.buffer, "This is a long string");
  b.length = strlen(b.buffer);

  const char *src2 = " that will fill the whole buffer and leave no space for "
                     "some of the chars.";
  int result2 = smart_append(&b, src2);

  printf("Example 2:\n");
  printf("Src 2: %s\n", src2);
  printf("  result: %d\n", result2);
  printf("  buffer: \"%s\"\n", b.buffer);
  printf("  length: %zu\n", b.length);

  return 0;
}
