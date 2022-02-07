package streams

func validateTopicName(validators ...func(topic, group string) error) func(topic, group string) error {
	return func(topic, group string) error {
		for _, v := range validators {
			if err := v(topic, group); err != nil {
				return err
			}
		}

		return nil
	}
}
