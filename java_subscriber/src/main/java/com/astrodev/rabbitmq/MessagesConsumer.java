package com.astrodev.rabbitmq;

import io.smallrye.mutiny.infrastructure.Infrastructure;
import jakarta.enterprise.context.ApplicationScoped;
import org.eclipse.microprofile.reactive.messaging.Acknowledgment;
import org.eclipse.microprofile.reactive.messaging.Incoming;
import org.eclipse.microprofile.reactive.messaging.Message;
import org.jboss.logging.Logger;
import io.smallrye.mutiny.Uni;

import java.time.Duration;
import java.util.concurrent.CompletionStage;

@ApplicationScoped
public class MessagesConsumer {
    private static final Logger logger = Logger.getLogger("messages");

    @Incoming("messages")
    public Uni<Void> consume(Message<String> message) {
        String payload = message.getPayload();
        long secondsToWait = payload.chars().filter(letter -> letter == '.').count();

        Uni<Void> work = Uni.createFrom().voidItem();
        if (secondsToWait > 0) {
            work = work.onItem().delayIt().by(Duration.ofSeconds(secondsToWait));
        }

        return work
                .onItem().invoke(() -> logger.info(payload))
                .onItem().transformToUni(ignored -> Uni.createFrom().completionStage(message::ack))
                .runSubscriptionOn(Infrastructure.getDefaultWorkerPool()); // Run on a worker thread
    }
}
